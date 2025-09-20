package handlers

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.com/m/v2/internal/domain"
	"example.com/m/v2/internal/parsers"
	"example.com/m/v2/internal/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader     websocket.Upgrader
	redisService *services.RedisService
	dataService  *services.DataService
	csvParser    *parsers.CSVParser
	connections  map[string]*Connection
	mu           sync.RWMutex
}

type Connection struct {
	ID       string
	Conn     *websocket.Conn
	Send     chan []byte
	LastSeen time.Time
}

type Message struct {
	Type      string          `json:"type"`
	SessionID string          `json:"session_id,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     string          `json:"error,omitempty"`
}

type UploadMessage struct {
	FileType string          `json:"file_type"`
	FileName string          `json:"file_name"`
	Chunks   []ChunkMessage  `json:"chunks"`
}

type ChunkMessage struct {
	Index int    `json:"index"`
	Data  string `json:"data"` // Base64 encoded CSV data
	Total int    `json:"total"`
}

type ProgressMessage struct {
	JobID          string `json:"job_id"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	Progress       int    `json:"progress"`
	TotalItems     int    `json:"total_items"`
	ProcessedItems int    `json:"processed_items"`
	Message        string `json:"message"`
}

type BatchConfig struct {
	Estados    int `json:"estados"`
	Municipios int `json:"municipios"`
	Medicos    int `json:"medicos"`
	Hospitais  int `json:"hospitais"`
	Pacientes  int `json:"pacientes"`
	CID10      int `json:"cid10"`
}

func NewWebSocketHandler(redisService *services.RedisService, dataService *services.DataService) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Em produção, configurar adequadamente
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		redisService: redisService,
		dataService:  dataService,
		csvParser:    parsers.NewCSVParser(),
		connections:  make(map[string]*Connection),
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	connectionID := uuid.New().String()
	connection := &Connection{
		ID:       connectionID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		LastSeen: time.Now(),
	}

	h.mu.Lock()
	h.connections[connectionID] = connection
	h.mu.Unlock()

	// Iniciar goroutines para leitura e escrita
	go h.writePump(connection)
	go h.readPump(connection)

	log.Printf("WebSocket connection established: %s", connectionID)
}

func (h *WebSocketHandler) readPump(conn *Connection) {
	defer func() {
		h.closeConnection(conn)
	}()

	conn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.Conn.SetPongHandler(func(string) error {
		conn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		err := conn.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		conn.LastSeen = time.Now()
		h.handleMessage(conn, &msg)
	}
}

func (h *WebSocketHandler) writePump(conn *Connection) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) handleMessage(conn *Connection, msg *Message) {
	switch msg.Type {
	case "upload_start":
		h.handleUploadStart(conn, msg)
	case "upload_chunk":
		h.handleUploadChunk(conn, msg)
	case "upload_complete":
		h.handleUploadComplete(conn, msg)
	case "get_progress":
		h.handleGetProgress(conn, msg)
	case "get_queue_status":
		h.handleGetQueueStatus(conn, msg)
	default:
		h.sendError(conn, "Unknown message type: "+msg.Type)
	}
}

func (h *WebSocketHandler) handleUploadStart(conn *Connection, msg *Message) {
	sessionID := msg.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	response := Message{
		Type:      "upload_ready",
		SessionID: sessionID,
	}

	h.sendMessage(conn, &response)
}

func (h *WebSocketHandler) handleUploadChunk(conn *Connection, msg *Message) {
	// Process chunk and send acknowledgment
	response := Message{
		Type:      "chunk_received",
		SessionID: msg.SessionID,
	}

	h.sendMessage(conn, &response)
}

func (h *WebSocketHandler) handleUploadComplete(conn *Connection, msg *Message) {
	var uploadMsg UploadMessage
	if err := json.Unmarshal(msg.Data, &uploadMsg); err != nil {
		h.sendError(conn, "Invalid upload data: "+err.Error())
		return
	}

	log.Printf("Upload recebido - Arquivo: %s, Tipo: %s, Chunks: %d", uploadMsg.FileName, uploadMsg.FileType, len(uploadMsg.Chunks))

	// Reconstituir dados do arquivo
	csvData, err := h.reconstructFile(uploadMsg.Chunks)
	if err != nil {
		h.sendError(conn, "Failed to reconstruct file: "+err.Error())
		return
	}

	// Processar baseado no tipo de arquivo
	log.Printf("Processando arquivo do tipo: %s", uploadMsg.FileType)
	err = h.processUploadedData(msg.SessionID, uploadMsg.FileType, csvData)
	if err != nil {
		h.sendError(conn, "Failed to process data: "+err.Error())
		return
	}

	response := Message{
		Type:      "upload_processing",
		SessionID: msg.SessionID,
		Data:      []byte(`{"message":"File uploaded and queued for processing"}`),
	}

	h.sendMessage(conn, &response)
}

func (h *WebSocketHandler) handleGetProgress(conn *Connection, msg *Message) {
	jobs, err := h.redisService.GetSessionJobs(msg.SessionID)
	if err != nil {
		h.sendError(conn, "Failed to get progress: "+err.Error())
		return
	}

	progressData := make([]ProgressMessage, len(jobs))
	for i, job := range jobs {
		progressData[i] = ProgressMessage{
			JobID:          job.ID,
			Type:           job.Type,
			Status:         string(job.Status),
			Progress:       job.Progress,
			TotalItems:     job.TotalItems,
			ProcessedItems: job.ProcessedItems,
			Message:        job.Error,
		}
	}

	data, _ := json.Marshal(progressData)
	response := Message{
		Type:      "progress_update",
		SessionID: msg.SessionID,
		Data:      data,
	}

	h.sendMessage(conn, &response)
}

func (h *WebSocketHandler) handleGetQueueStatus(conn *Connection, msg *Message) {
	queueStats := make(map[string]int64)
	jobTypes := []string{"estados", "municipios", "medicos", "hospitais", "pacientes", "cid10"}

	for _, jobType := range jobTypes {
		length, _ := h.redisService.GetQueueLength(jobType)
		queueStats[jobType] = length
	}

	data, _ := json.Marshal(queueStats)
	response := Message{
		Type: "queue_status",
		Data: data,
	}

	h.sendMessage(conn, &response)
}

func (h *WebSocketHandler) reconstructFile(chunks []ChunkMessage) (string, error) {
	if len(chunks) == 0 {
		return "", fmt.Errorf("no chunks provided")
	}

	log.Printf("Reconstruindo arquivo de %d chunks", len(chunks))

	// Ordenar chunks por índice
	sortedChunks := make([]string, chunks[0].Total)
	for _, chunk := range chunks {
		if chunk.Index >= len(sortedChunks) {
			return "", fmt.Errorf("invalid chunk index: %d", chunk.Index)
		}
		sortedChunks[chunk.Index] = chunk.Data
	}

	// Concatenar dados base64
	var base64Data string
	for _, data := range sortedChunks {
		base64Data += data
	}

	log.Printf("Tamanho do base64 concatenado: %d bytes", len(base64Data))

	// Decodificar base64 para obter o conteúdo real do CSV
	csvBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	csvContent := string(csvBytes)
	log.Printf("Arquivo CSV decodificado com sucesso, tamanho: %d bytes", len(csvContent))
	log.Printf("Primeiras 200 caracteres do CSV: %s", csvContent[:min(200, len(csvContent))])

	return csvContent, nil
}

// Helper function para min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *WebSocketHandler) processUploadedData(sessionID, fileType, csvData string) error {
	// Configuração de batch por tipo
	batchSizes := map[string]int{
		"estados":    50,
		"municipios": 200,
		"medicos":    300,
		"hospitais":  150,
		"pacientes":  250,
		"cid10":      400,
	}

	batchSize := batchSizes[fileType]
	if batchSize == 0 {
		batchSize = 100 // Default
	}

	// Processar baseado no tipo de arquivo
	log.Printf("Switch case - tipo recebido: '%s' (length: %d)", fileType, len(fileType))

	// Normalizar o tipo para lowercase e remover espaços
	normalizedType := strings.ToLower(strings.TrimSpace(fileType))
	log.Printf("Tipo normalizado: '%s'", normalizedType)

	switch normalizedType {
	case "municipios":
		log.Printf("Processando como MUNICIPIOS")
		return h.processMunicipios(sessionID, csvData, batchSize)
	case "medicos":
		log.Printf("Processando como MEDICOS")
		return h.processMedicos(sessionID, csvData, batchSize)
	case "estados":
		log.Printf("Processando como ESTADOS")
		return h.processEstados(sessionID, csvData, batchSize)
	case "hospitais":
		log.Printf("Processando como HOSPITAIS")
		return h.processHospitais(sessionID, csvData, batchSize)
	case "pacientes":
		log.Printf("Processando como PACIENTES")
		return h.processPacientes(sessionID, csvData, batchSize)
	case "cid10":
		log.Printf("Processando como CID10")
		return h.processCID10(sessionID, csvData, batchSize)
	default:
		log.Printf("ERRO: Tipo não reconhecido: '%s'", fileType)
		return fmt.Errorf("unsupported file type: %s", fileType)
	}
}

func (h *WebSocketHandler) processMunicipios(sessionID, csvData string, batchSize int) error {
	// Parse CSV data usando o parser existente
	// Primeiro, criar um arquivo temporário em memória
	municipios, err := h.parseCSVMunicipiosFromString(csvData)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(municipios) == 0 {
		return fmt.Errorf("no valid municipalities found in CSV")
	}

	totalItems := len(municipios)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processing %d municipalities in %d batches for session %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		// Extrair batch atual
		batchMunicipios := municipios[start:end]

		// Criar job para o batch
		jobID := uuid.New().String()
		job := &services.Job{
			ID:           jobID,
			Type:         "municipios",
			Status:       services.JobStatusPending,
			TotalItems:   len(batchMunicipios),
			CreatedAt:    time.Now(),
			SessionID:    sessionID,
		}

		// Serializar dados do batch
		batchDataBytes, err := json.Marshal(batchMunicipios)
		if err != nil {
			return fmt.Errorf("failed to marshal batch data: %v", err)
		}

		batchData := services.BatchData{
			BatchNumber:  i + 1,
			TotalBatches: totalBatches,
			Data:         batchDataBytes,
		}

		data, _ := json.Marshal(batchData)
		job.Data = data

		err = h.redisService.EnqueueJob(job)
		if err != nil {
			return fmt.Errorf("failed to enqueue job: %v", err)
		}

		log.Printf("Enqueued batch %d/%d with %d municipalities", i+1, totalBatches, len(batchMunicipios))
	}

	return nil
}

func (h *WebSocketHandler) processMedicos(sessionID, csvData string, batchSize int) error {
	medicos, err := h.parseCSVMedicosFromString(csvData)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(medicos) == 0 {
		return fmt.Errorf("no valid doctors found in CSV")
	}

	totalItems := len(medicos)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processing %d doctors in %d batches for session %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchMedicos := medicos[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:           jobID,
			Type:         "medicos",
			Status:       services.JobStatusPending,
			TotalItems:   len(batchMedicos),
			CreatedAt:    time.Now(),
			SessionID:    sessionID,
		}

		batchDataBytes, err := json.Marshal(batchMedicos)
		if err != nil {
			return fmt.Errorf("failed to marshal batch data: %v", err)
		}

		batchData := services.BatchData{
			BatchNumber:  i + 1,
			TotalBatches: totalBatches,
			Data:         batchDataBytes,
		}

		data, _ := json.Marshal(batchData)
		job.Data = data

		err = h.redisService.EnqueueJob(job)
		if err != nil {
			return fmt.Errorf("failed to enqueue job: %v", err)
		}

		log.Printf("Enqueued batch %d/%d with %d doctors", i+1, totalBatches, len(batchMedicos))
	}

	return nil
}

// parseCSVMunicipiosFromString parses CSV data from string
func (h *WebSocketHandler) parseCSVMunicipiosFromString(csvData string) ([]domain.Municipio, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ','
	reader.LazyQuotes = true

	// Ler cabeçalho
	headers, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("empty CSV data")
	}

	// Mapear índices das colunas
	headerMap := make(map[string]int)
	for i, header := range headers[0] {
		headerMap[strings.TrimSpace(header)] = i
	}

	var municipios []domain.Municipio

	// Processar linhas de dados (pular cabeçalho)
	for i := 1; i < len(headers); i++ {
		record := headers[i]

		municipio := domain.Municipio{}

		// Mapear campos usando os índices do cabeçalho
		if idx, exists := headerMap["codigo_ibge"]; exists && idx < len(record) {
			municipio.Codigo = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["nome"]; exists && idx < len(record) {
			municipio.Nome = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["latitude"]; exists && idx < len(record) {
			municipio.Latitude = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["longitude"]; exists && idx < len(record) {
			municipio.Longitude = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["capital"]; exists && idx < len(record) {
			municipio.Capital = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["codigo_uf"]; exists && idx < len(record) {
			municipio.CodigoUF = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["siafi_id"]; exists && idx < len(record) {
			municipio.SiafiId = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["ddd"]; exists && idx < len(record) {
			municipio.DDD = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["fuso_horario"]; exists && idx < len(record) {
			municipio.FusoHora = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["populacao"]; exists && idx < len(record) {
			if pop, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
				municipio.Populacao = pop
			}
		}

		// Só adicionar se tiver pelo menos código e nome
		if municipio.Codigo != "" && municipio.Nome != "" {
			municipios = append(municipios, municipio)
		}
	}

	return municipios, nil
}

// parseCSVMedicosFromString parses CSV data from string
func (h *WebSocketHandler) parseCSVMedicosFromString(csvData string) ([]domain.Medico, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ','
	reader.LazyQuotes = true

	headers, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("empty CSV data")
	}

	log.Printf("CSV headers para médicos: %v", headers[0])

	headerMap := make(map[string]int)
	for i, header := range headers[0] {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var medicos []domain.Medico

	for i := 1; i < len(headers); i++ {
		record := headers[i]
		medico := domain.Medico{}

		// Múltiplas variações para nome
		nomeColumns := []string{"nome_completo", "nome", "name", "medico_nome", "doctor_name", "full_name"}
		for _, col := range nomeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				medico.Nome = strings.TrimSpace(record[idx])
				break
			}
		}

		// Múltiplas variações para especialidade
		especialidadeColumns := []string{"especialidade", "specialty", "specialism", "area", "categoria"}
		for _, col := range especialidadeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				medico.Especialidade = strings.TrimSpace(record[idx])
				break
			}
		}

		// Múltiplas variações para código/cidade (usar cidade como fallback para município)
		municipioColumns := []string{"codigo", "cod_municipio", "codigo_municipio", "municipio_id", "cidade", "city", "city_code", "municipality_code"}
		for _, col := range municipioColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				medico.CodMunicipio = strings.TrimSpace(record[idx])
				break
			}
		}

		// Log detalhado do parsing
		log.Printf("Registro %d - Nome: '%s', Especialidade: '%s', CodMunicipio: '%s'", i, medico.Nome, medico.Especialidade, medico.CodMunicipio)

		// Só adicionar se tiver pelo menos nome
		if medico.Nome != "" {
			medico.UUID = uuid.New()
			medicos = append(medicos, medico)
			log.Printf("✓ Médico %d adicionado: %s", i, medico.Nome)
		} else {
			log.Printf("✗ Médico %d rejeitado - nome vazio", i)
		}
	}

	log.Printf("Total de médicos parseados: %d", len(medicos))
	return medicos, nil
}

// processEstados processa dados de estados
func (h *WebSocketHandler) processEstados(sessionID, csvData string, batchSize int) error {
	estados, err := h.parseCSVEstadosFromString(csvData)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(estados) == 0 {
		return fmt.Errorf("no valid states found in CSV")
	}

	totalItems := len(estados)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processing %d states in %d batches for session %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchEstados := estados[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:           jobID,
			Type:         "estados",
			Status:       services.JobStatusPending,
			TotalItems:   len(batchEstados),
			CreatedAt:    time.Now(),
			SessionID:    sessionID,
		}

		batchDataBytes, err := json.Marshal(batchEstados)
		if err != nil {
			return fmt.Errorf("failed to marshal batch data: %v", err)
		}

		batchData := services.BatchData{
			BatchNumber:  i + 1,
			TotalBatches: totalBatches,
			Data:         batchDataBytes,
		}

		data, _ := json.Marshal(batchData)
		job.Data = data

		err = h.redisService.EnqueueJob(job)
		if err != nil {
			return fmt.Errorf("failed to enqueue job: %v", err)
		}

		log.Printf("Enqueued batch %d/%d with %d states", i+1, totalBatches, len(batchEstados))
	}

	return nil
}

// parseCSVEstadosFromString parses CSV data from string
func (h *WebSocketHandler) parseCSVEstadosFromString(csvData string) ([]domain.Estado, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ','
	reader.LazyQuotes = true

	headers, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("empty CSV data")
	}

	// Log dos cabeçalhos para debug
	log.Printf("CSV headers encontrados: %v", headers[0])

	headerMap := make(map[string]int)
	for i, header := range headers[0] {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
		log.Printf("Header %d: '%s' -> '%s'", i, header, cleanHeader)
	}

	var estados []domain.Estado

	for i := 1; i < len(headers); i++ {
		record := headers[i]
		estado := domain.Estado{}

		// Tentar múltiplas variações dos nomes de colunas
		codigoColumns := []string{"codigo", "cod", "code", "id", "estado_id", "uf"}
		for _, col := range codigoColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				estado.Codigo = strings.TrimSpace(record[idx])
				break
			}
		}

		ufColumns := []string{"unidade_federativa", "uf", "sigla", "estado", "state"}
		for _, col := range ufColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				estado.UnidadeFederativa = strings.TrimSpace(record[idx])
				break
			}
		}

		nomeColumns := []string{"nome", "name", "estado_nome", "estado", "state_name"}
		for _, col := range nomeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				estado.Nome = strings.TrimSpace(record[idx])
				break
			}
		}

		regiaoColumns := []string{"regiao", "region", "macroregiao", "macro_regiao"}
		for _, col := range regiaoColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				estado.Regiao = strings.TrimSpace(record[idx])
				break
			}
		}

		latColumns := []string{"latitude", "lat", "latitude_decimal"}
		for _, col := range latColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				estado.Latitude = strings.TrimSpace(record[idx])
				break
			}
		}

		lngColumns := []string{"longitude", "lng", "lon", "longitude_decimal"}
		for _, col := range lngColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				estado.Longitude = strings.TrimSpace(record[idx])
				break
			}
		}

		// Log do registro parseado para debug
		log.Printf("Registro %d parseado: Codigo='%s', Nome='%s', UF='%s'", i, estado.Codigo, estado.Nome, estado.UnidadeFederativa)

		// Adicionar se tiver pelo menos um identificador (código, UF ou nome)
		if estado.Codigo != "" || estado.UnidadeFederativa != "" || estado.Nome != "" {
			estados = append(estados, estado)
		}
	}

	log.Printf("Total de estados parseados: %d", len(estados))
	return estados, nil
}

// Placeholder implementations for other types
func (h *WebSocketHandler) processHospitais(sessionID, csvData string, batchSize int) error {
	// TODO: Implementar quando necessário
	return fmt.Errorf("hospital processing not implemented yet")
}

func (h *WebSocketHandler) processPacientes(sessionID, csvData string, batchSize int) error {
	// TODO: Implementar quando necessário
	return fmt.Errorf("patient processing not implemented yet")
}

func (h *WebSocketHandler) processCID10(sessionID, csvData string, batchSize int) error {
	// TODO: Implementar quando necessário
	return fmt.Errorf("CID10 processing not implemented yet")
}

func (h *WebSocketHandler) sendMessage(conn *Connection, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case conn.Send <- data:
	default:
		h.closeConnection(conn)
	}
}

func (h *WebSocketHandler) sendError(conn *Connection, error string) {
	msg := Message{
		Type:  "error",
		Error: error,
	}
	h.sendMessage(conn, &msg)
}

func (h *WebSocketHandler) closeConnection(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.connections[conn.ID]; ok {
		delete(h.connections, conn.ID)
		close(conn.Send)
		conn.Conn.Close()
		log.Printf("WebSocket connection closed: %s", conn.ID)
	}
}

// BroadcastProgress envia atualizações de progresso para todas as conexões
func (h *WebSocketHandler) BroadcastProgress(sessionID string, progress ProgressMessage) {
	msg := Message{
		Type:      "progress_update",
		SessionID: sessionID,
	}

	data, _ := json.Marshal([]ProgressMessage{progress})
	msg.Data = data

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.connections {
		h.sendMessage(conn, &msg)
	}
}