package handlers

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// getEnvAsInt retrieves an environment variable as an integer, with a default fallback
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

type WebSocketHandler struct {
	upgrader      websocket.Upgrader
	redisService  *services.RedisService
	dataService   *services.DataService
	csvParser     *parsers.CSVParser
	unifiedParser *parsers.UnifiedParser
	connections   map[string]*Connection
	mu            sync.RWMutex
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
	FileType string         `json:"file_type"`
	FileName string         `json:"file_name"`
	Chunks   []ChunkMessage `json:"chunks"`
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
		redisService:  redisService,
		dataService:   dataService,
		csvParser:     parsers.NewCSVParser(),
		unifiedParser: parsers.NewUnifiedParser(),
		connections:   make(map[string]*Connection),
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

	log.Printf("🔄 Reconstruindo arquivo de %d chunks", len(chunks))

	// Ordenar chunks por índice usando mapa para melhor performance
	chunkMap := make(map[int]string, len(chunks))
	maxIndex := 0

	for _, chunk := range chunks {
		chunkMap[chunk.Index] = chunk.Data
		if chunk.Index > maxIndex {
			maxIndex = chunk.Index
		}
	}

	// Usar strings.Builder para concatenação mais eficiente
	var builder strings.Builder
	totalSize := 0
	for _, chunk := range chunks {
		totalSize += len(chunk.Data)
	}
	builder.Grow(totalSize) // Pre-alocar tamanho necessário

	// Concatenar chunks na ordem correta
	for i := 0; i <= maxIndex; i++ {
		if data, exists := chunkMap[i]; exists {
			builder.WriteString(data)
		} else {
			return "", fmt.Errorf("chunk %d está faltando", i)
		}
	}

	base64Data := builder.String()
	log.Printf("📊 Base64 concatenado: %d bytes", len(base64Data))

	// Decodificar base64 com buffer pre-alocado
	expectedSize := (len(base64Data) * 3) / 4
	csvBytes := make([]byte, expectedSize)

	n, err := base64.StdEncoding.Decode(csvBytes, []byte(base64Data))
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar base64: %v", err)
	}

	csvContent := string(csvBytes[:n])
	log.Printf("✅ Arquivo decodificado: %d bytes (%d chunks processados)", len(csvContent), len(chunks))

	// Log da estrutura para debug (apenas primeiras linhas)
	lines := strings.Split(csvContent, "\n")
	if len(lines) > 0 {
		log.Printf("📋 Primeira linha: %s", lines[0][:min(100, len(lines[0]))])
	}

	return csvContent, nil
}


func (h *WebSocketHandler) processUploadedData(sessionID, fileType, csvData string) error {
	// Tentar usar o parser unificado primeiro
	data := []byte(csvData)
	filename := fmt.Sprintf("%s.xml", fileType) // Usar .xml baseado no que vemos nos logs

	log.Printf("🔍 Debug: tentando detectar formato para %s (tamanho: %d bytes)", filename, len(data))
	log.Printf("🔍 Debug: primeiros 100 chars: %s", string(data[:min(100, len(data))]))

	parsedData, detectedFormat, err := h.unifiedParser.ParseFile(data, filename)
	if err == nil {
		log.Printf("✅ Formato detectado automaticamente: %s", detectedFormat.String())
		return h.processUnifiedData(sessionID, fileType, parsedData, detectedFormat)
	}

	log.Printf("⚠ Parser unificado falhou (%v), usando parser legado", err)
	
	// Fallback para o processamento legado
	return h.processUploadedDataLegacy(sessionID, fileType, csvData)
}

// processUnifiedData processa dados usando o parser unificado
func (h *WebSocketHandler) processUnifiedData(sessionID, fileType string, parsedData *parsers.ParsedData, format parsers.FileFormat) error {
	log.Printf("Processando dados unificados - Tipo: %s, Formato: %s, Linhas: %d, Colunas: %d", 
		fileType, format.String(), len(parsedData.Rows), len(parsedData.Headers))

	// Log dos headers detectados
	log.Printf("Headers detectados: %v", parsedData.Headers)

	// Configuração de batch por tipo
	batchSizes := map[string]int{
		"estados":    getEnvAsInt("BATCH_SIZE_ESTADOS", 50),
		"municipios": getEnvAsInt("BATCH_SIZE_MUNICIPIOS", 200),
		"medicos":    getEnvAsInt("BATCH_SIZE_MEDICOS", 300),
		"hospitais":  getEnvAsInt("BATCH_SIZE_HOSPITAIS", 150),
		"pacientes":  getEnvAsInt("BATCH_SIZE_PACIENTES", 250),
		"cid10":      getEnvAsInt("BATCH_SIZE_CID10", 400),
	}

	batchSize := batchSizes[fileType]
	if batchSize == 0 {
		batchSize = 100 // Default
	}

	// Processar baseado no tipo normalizado
	normalizedType := strings.ToLower(strings.TrimSpace(fileType))
	
	// Para FHIR, primeiro tentar mapear automaticamente
	if format == parsers.FormatFHIRJSON || format == parsers.FormatFHIRXML {
		return h.processFHIRData(sessionID, parsedData, batchSize)
	}

	// Para HL7, processar automaticamente
	if format == parsers.FormatHL7 {
		return h.processHL7Data(sessionID, parsedData, batchSize)
	}

	// Para outros formatos, usar conversão baseada no tipo
	switch normalizedType {
	case "hospitais", "hospitals":
		return h.processUnifiedHospitals(sessionID, parsedData, batchSize)
	case "pacientes", "patients":
		return h.processUnifiedPatients(sessionID, parsedData, batchSize)
	case "medicos", "doctors", "practitioners":
		return h.processUnifiedMedicos(sessionID, parsedData, batchSize)
	case "municipios", "municipalities":
		return h.processUnifiedMunicipios(sessionID, parsedData, batchSize)
	case "estados", "states":
		return h.processUnifiedEstados(sessionID, parsedData, batchSize)
	case "cid10":
		return h.processUnifiedCID10(sessionID, parsedData, batchSize)
	default:
		return fmt.Errorf("tipo não suportado pelo parser unificado: %s", fileType)
	}
}

// processUnifiedHospitals processa hospitais usando dados unificados
func (h *WebSocketHandler) processUnifiedHospitals(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	hospitais, err := h.convertRowsToHospitals(parsedData)
	if err != nil {
		return fmt.Errorf("erro ao converter dados para hospitais: %v", err)
	}

	if len(hospitais) == 0 {
		return fmt.Errorf("nenhum hospital válido encontrado nos dados")
	}

	totalItems := len(hospitais)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processando %d hospitais em %d batches para sessão %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchHospitais := hospitais[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "hospitais",
			Status:     services.JobStatusPending,
			TotalItems: len(batchHospitais),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchHospitais)
		if err != nil {
			return fmt.Errorf("erro ao serializar batch: %v", err)
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
			return fmt.Errorf("erro ao enfileirar job: %v", err)
		}

		log.Printf("✓ Batch %d/%d enfileirado com %d hospitais", i+1, totalBatches, len(batchHospitais))
	}

	return nil
}

// convertRowsToHospitals converte dados tabulares para hospitais
func (h *WebSocketHandler) convertRowsToHospitals(parsedData *parsers.ParsedData) ([]domain.Hospital, error) {
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var hospitais []domain.Hospital

	for rowIndex, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}

		hospital := domain.Hospital{}

		// UUID/Código
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			uuidStr := strings.TrimSpace(row[idx])
			if uuidStr != "" {
				if parsedUUID, err := uuid.Parse(uuidStr); err == nil {
					hospital.UUID = parsedUUID
				} else {
					hospital.UUID = uuid.New()
				}
			} else {
				hospital.UUID = uuid.New()
			}
		} else {
			hospital.UUID = uuid.New()
		}

		// Nome
		nomeFields := []string{"nome", "hospital_nome", "hospital", "name", "hospital_name"}
		for _, field := range nomeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				hospital.Nome = strings.TrimSpace(row[idx])
				break
			}
		}

		// CEP
		cepFields := []string{"cep", "zipcode", "postal_code"}
		for _, field := range cepFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				hospital.CEP = strings.TrimSpace(row[idx])
				break
			}
		}

		// Especialidades
		especialidadeFields := []string{"especialidades", "specialties", "services", "servicos"}
		for _, field := range especialidadeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				hospital.Especialidades = strings.TrimSpace(row[idx])
				break
			}
		}

		// Leitos totais
		leitosFields := []string{"leitos_totais", "leitos", "beds", "total_beds"}
		for _, field := range leitosFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				if leitos, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
					hospital.LeitosTotais = leitos
				}
				break
			}
		}

		// Código do município
		municipioFields := []string{"cidade", "cod_municipio", "codigo_municipio", "municipio_id", "city", "city_code"}
		for _, field := range municipioFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				hospital.CodMunicipio = strings.TrimSpace(row[idx])
				break
			}
		}

		// Bairro
		bairroFields := []string{"bairro", "district", "neighborhood"}
		for _, field := range bairroFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				hospital.Bairro = strings.TrimSpace(row[idx])
				break
			}
		}

		// Só adicionar se tiver pelo menos nome
		if hospital.Nome != "" {
			hospitais = append(hospitais, hospital)
		} else {
			log.Printf("⚠ Linha %d ignorada - hospital sem nome", rowIndex+1)
		}
	}

	return hospitais, nil
}

// processUnifiedPatients processa pacientes usando dados unificados
func (h *WebSocketHandler) processUnifiedPatients(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	// Verificar se é XML estruturado e usar converter específico
	if detectedType, exists := parsedData.Metadata["detected_type"]; exists && detectedType == "pacientes" {
		return h.processUnifiedPatientsFromXML(sessionID, parsedData, batchSize)
	}
	
	// Fallback para conversão genérica
	pacientes, err := h.convertRowsToPatients(parsedData)
	if err != nil {
		return fmt.Errorf("erro ao converter dados para pacientes: %v", err)
	}

	if len(pacientes) == 0 {
		return fmt.Errorf("nenhum paciente válido encontrado nos dados")
	}

	return h.enqueuePatientsJobs(sessionID, pacientes, batchSize)
}

// processUnifiedPatientsFromXML processa pacientes usando converter XML específico
func (h *WebSocketHandler) processUnifiedPatientsFromXML(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	// Usar o converter XML específico
	xmlParser := parsers.NewXMLParser()
	pacientes, err := xmlParser.ConvertXMLToDomain(parsedData, "pacientes")
	if err != nil {
		return fmt.Errorf("erro ao converter XML para pacientes: %v", err)
	}

	pacientesList, ok := pacientes.([]domain.Paciente)
	if !ok {
		return fmt.Errorf("erro de tipo na conversão de pacientes")
	}

	if len(pacientesList) == 0 {
		return fmt.Errorf("nenhum paciente válido encontrado no XML")
	}

	log.Printf("✓ Convertidos %d pacientes do XML estruturado", len(pacientesList))
	
	return h.enqueuePatientsJobs(sessionID, pacientesList, batchSize)
}

// enqueuePatientsJobs enfileira jobs de pacientes
func (h *WebSocketHandler) enqueuePatientsJobs(sessionID string, pacientes []domain.Paciente, batchSize int) error {
	totalItems := len(pacientes)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processando %d pacientes em %d batches para sessão %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchPacientes := pacientes[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "pacientes",
			Status:     services.JobStatusPending,
			TotalItems: len(batchPacientes),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchPacientes)
		if err != nil {
			return fmt.Errorf("erro ao serializar batch: %v", err)
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
			return fmt.Errorf("erro ao enfileirar job: %v", err)
		}

		log.Printf("✓ Batch %d/%d enfileirado com %d pacientes", i+1, totalBatches, len(batchPacientes))
	}

	return nil
}

// convertRowsToPatients converte dados tabulares genéricos para pacientes
func (h *WebSocketHandler) convertRowsToPatients(parsedData *parsers.ParsedData) ([]domain.Paciente, error) {
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var pacientes []domain.Paciente

	for rowIndex, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}

		paciente := domain.Paciente{}

		// ID/Código (UUID)
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			uuidStr := strings.TrimSpace(row[idx])
			if uuidStr != "" {
				if parsedUUID, err := uuid.Parse(uuidStr); err == nil {
					paciente.ID = parsedUUID
				} else {
					paciente.ID = uuid.New()
				}
			} else {
				paciente.ID = uuid.New()
			}
		} else {
			paciente.ID = uuid.New()
		}

		// CPF
		cpfFields := []string{"cpf", "documento", "doc"}
		for _, field := range cpfFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.CPF = strings.TrimSpace(row[idx])
				break
			}
		}

		// Nome
		nomeFields := []string{"nome", "nome_completo", "name", "full_name", "patient_name"}
		for _, field := range nomeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.Nome = strings.TrimSpace(row[idx])
				break
			}
		}

		// Gênero
		generoFields := []string{"genero", "gender", "sex", "sexo"}
		for _, field := range generoFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.Genero = strings.TrimSpace(row[idx])
				break
			}
		}

		// Código do município
		municipioFields := []string{"cod_municipio", "codigo_municipio", "municipio_id", "city_code"}
		for _, field := range municipioFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.CodMunicipio = strings.TrimSpace(row[idx])
				break
			}
		}

		// Bairro
		bairroFields := []string{"bairro", "district", "neighborhood"}
		for _, field := range bairroFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.Bairro = strings.TrimSpace(row[idx])
				break
			}
		}

		// Convênio
		convenioFields := []string{"convenio", "insurance", "plano_saude"}
		for _, field := range convenioFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.Convenio = strings.TrimSpace(row[idx])
				break
			}
		}

		// CID10
		cid10Fields := []string{"cid10", "cid-10", "diagnosis", "diagnostico"}
		for _, field := range cid10Fields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				paciente.CID10 = strings.TrimSpace(row[idx])
				break
			}
		}

		// Só adicionar se tiver pelo menos CPF e nome
		if paciente.CPF != "" && paciente.Nome != "" {
			pacientes = append(pacientes, paciente)
		} else {
			log.Printf("⚠ Linha %d ignorada - paciente sem CPF ou nome", rowIndex+1)
		}
	}

	return pacientes, nil
}

// convertRowsToMedicos converte dados tabulares para médicos
func (h *WebSocketHandler) convertRowsToMedicos(parsedData *parsers.ParsedData) ([]domain.Medico, error) {
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var medicos []domain.Medico

	for rowIndex, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}

		medico := domain.Medico{}

		// UUID/Código
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			uuidStr := strings.TrimSpace(row[idx])
			if uuidStr != "" {
				if parsedUUID, err := uuid.Parse(uuidStr); err == nil {
					medico.UUID = parsedUUID
				} else {
					medico.UUID = uuid.New()
				}
			} else {
				medico.UUID = uuid.New()
			}
		} else {
			medico.UUID = uuid.New()
		}

		// Nome
		nomeFields := []string{"nome", "nome_completo", "name", "full_name", "medico_nome", "doctor_name"}
		for _, field := range nomeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				medico.Nome = strings.TrimSpace(row[idx])
				break
			}
		}

		// Especialidade
		especialidadeFields := []string{"especialidade", "specialty", "specialization", "area", "categoria"}
		for _, field := range especialidadeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				medico.Especialidade = strings.TrimSpace(row[idx])
				break
			}
		}

		// Código do município
		municipioFields := []string{"cidade", "cod_municipio", "codigo_municipio", "municipio_id", "city", "city_code"}
		for _, field := range municipioFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				medico.CodMunicipio = strings.TrimSpace(row[idx])
				break
			}
		}

		// Só adicionar se tiver pelo menos nome
		if medico.Nome != "" {
			medicos = append(medicos, medico)
		} else {
			log.Printf("⚠ Linha %d ignorada - médico sem nome", rowIndex+1)
		}
	}

	return medicos, nil
}

// convertRowsToMunicipios converte dados tabulares para municípios
func (h *WebSocketHandler) convertRowsToMunicipios(parsedData *parsers.ParsedData) ([]domain.Municipio, error) {
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var municipios []domain.Municipio

	for rowIndex, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}

		municipio := domain.Municipio{}

		// Código IBGE
		codigoFields := []string{"codigo_ibge", "codigo", "code", "municipio_id", "ibge_code"}
		for _, field := range codigoFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.Codigo = strings.TrimSpace(row[idx])
				break
			}
		}

		// Nome
		nomeFields := []string{"nome", "name", "municipio_nome", "municipality_name"}
		for _, field := range nomeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.Nome = strings.TrimSpace(row[idx])
				break
			}
		}

		// Latitude
		latFields := []string{"latitude", "lat", "latitude_decimal"}
		for _, field := range latFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.Latitude = strings.TrimSpace(row[idx])
				break
			}
		}

		// Longitude
		lngFields := []string{"longitude", "lng", "lon", "longitude_decimal"}
		for _, field := range lngFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.Longitude = strings.TrimSpace(row[idx])
				break
			}
		}

		// Capital
		capitalFields := []string{"capital", "is_capital", "e_capital"}
		for _, field := range capitalFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.Capital = strings.TrimSpace(row[idx])
				break
			}
		}

		// Código UF
		ufFields := []string{"codigo_uf", "uf_code", "state_code", "estado_codigo"}
		for _, field := range ufFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.CodigoUF = strings.TrimSpace(row[idx])
				break
			}
		}

		// SIAFI ID
		siafiFields := []string{"siafi_id", "siafi", "codigo_siafi"}
		for _, field := range siafiFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.SiafiId = strings.TrimSpace(row[idx])
				break
			}
		}

		// DDD
		dddFields := []string{"ddd", "area_code", "codigo_area"}
		for _, field := range dddFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.DDD = strings.TrimSpace(row[idx])
				break
			}
		}

		// Fuso horário
		fusoFields := []string{"fuso_horario", "fuso_hora", "timezone", "time_zone"}
		for _, field := range fusoFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				municipio.FusoHora = strings.TrimSpace(row[idx])
				break
			}
		}

		// População
		popFields := []string{"populacao", "population", "habitantes"}
		for _, field := range popFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				if pop, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
					municipio.Populacao = pop
				}
				break
			}
		}

		// Só adicionar se tiver código e nome
		if municipio.Codigo != "" && municipio.Nome != "" {
			municipios = append(municipios, municipio)
		} else {
			log.Printf("⚠ Linha %d ignorada - município sem código ou nome", rowIndex+1)
		}
	}

	return municipios, nil
}

// convertRowsToEstados converte dados tabulares para estados
func (h *WebSocketHandler) convertRowsToEstados(parsedData *parsers.ParsedData) ([]domain.Estado, error) {
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	log.Printf("🔍 DEBUG Estados - Headers disponíveis: %v", parsedData.Headers)
	log.Printf("🔍 DEBUG Estados - Header map: %v", headerMap)

	var estados []domain.Estado

	for rowIndex, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}

		estado := domain.Estado{}

		// Log da linha atual para debug
		log.Printf("🔍 DEBUG Estados - Linha %d: %v", rowIndex+1, row)

		// Código
		codigoFields := []string{"codigo", "codigo_uf", "code", "state_code", "uf_code"}
		for _, field := range codigoFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				codigoValue := strings.TrimSpace(row[idx])
				// Garantir que o código tenha no máximo 2 caracteres
				if len(codigoValue) > 2 {
					codigoValue = codigoValue[:2]
				}
				estado.Codigo = codigoValue
				log.Printf("🔍 DEBUG Estados - Linha %d campo '%s' (idx %d): '%s' -> codigo='%s'", rowIndex+1, field, idx, row[idx], estado.Codigo)
				break
			}
		}

		// Unidade Federativa (sigla)
		ufFields := []string{"unidade_federativa", "uf", "sigla", "state_abbr", "abbreviation"}
		for _, field := range ufFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				ufValue := strings.TrimSpace(row[idx])
				// Garantir que a UF tenha no máximo 2 caracteres
				if len(ufValue) > 2 {
					ufValue = ufValue[:2]
				}
				estado.UnidadeFederativa = ufValue
				log.Printf("🔍 DEBUG Estados - Linha %d campo '%s' (idx %d): '%s' -> uf='%s'", rowIndex+1, field, idx, row[idx], estado.UnidadeFederativa)
				break
			}
		}

		// Nome
		nomeFields := []string{"nome", "name", "estado_nome", "state_name"}
		for _, field := range nomeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				estado.Nome = strings.TrimSpace(row[idx])
				log.Printf("🔍 DEBUG Estados - Linha %d campo '%s' (idx %d): '%s' -> nome='%s'", rowIndex+1, field, idx, row[idx], estado.Nome)
				break
			}
		}

		// Região
		regiaoFields := []string{"regiao", "region", "macroregiao", "macro_regiao"}
		for _, field := range regiaoFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				estado.Regiao = strings.TrimSpace(row[idx])
				log.Printf("🔍 DEBUG Estados - Linha %d campo '%s' (idx %d): '%s' -> regiao='%s'", rowIndex+1, field, idx, row[idx], estado.Regiao)
				break
			}
		}

		// Latitude
		latFields := []string{"latitude", "lat", "latitude_decimal"}
		for _, field := range latFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				estado.Latitude = strings.TrimSpace(row[idx])
				break
			}
		}

		// Longitude
		lngFields := []string{"longitude", "lng", "lon", "longitude_decimal"}
		for _, field := range lngFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				estado.Longitude = strings.TrimSpace(row[idx])
				break
			}
		}

		// Log do estado final parseado
		log.Printf("✅ DEBUG Estados - Linha %d parseada: Codigo='%s', UF='%s', Nome='%s', Regiao='%s'",
			rowIndex+1, estado.Codigo, estado.UnidadeFederativa, estado.Nome, estado.Regiao)

		// Garantir que campos obrigatórios sejam preenchidos
		// Se não tiver código, usar UF como código (padrão comum)
		if estado.Codigo == "" && estado.UnidadeFederativa != "" {
			estado.Codigo = estado.UnidadeFederativa
			log.Printf("🔧 DEBUG Estados - Linha %d: usando UF como código: '%s'", rowIndex+1, estado.Codigo)
		}

		// Se ainda não tiver região, usar um padrão
		if estado.Regiao == "" {
			estado.Regiao = "Indefinida"
		}

		// Só adicionar se tiver pelo menos código e nome
		if estado.Codigo != "" && estado.Nome != "" {
			estados = append(estados, estado)
			log.Printf("✅ DEBUG Estados - Linha %d ADICIONADA ao slice", rowIndex+1)
		} else {
			log.Printf("⚠ Linha %d ignorada - estado sem código ou nome válido (codigo='%s', nome='%s')",
				rowIndex+1, estado.Codigo, estado.Nome)
		}
	}

	log.Printf("📊 DEBUG Estados - Total convertido: %d estados", len(estados))
	return estados, nil
}

func (h *WebSocketHandler) processUnifiedMedicos(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	medicos, err := h.convertRowsToMedicos(parsedData)
	if err != nil {
		return fmt.Errorf("erro ao converter dados para médicos: %v", err)
	}

	if len(medicos) == 0 {
		return fmt.Errorf("nenhum médico válido encontrado nos dados")
	}

	totalItems := len(medicos)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processando %d médicos em %d batches para sessão %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchMedicos := medicos[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "medicos",
			Status:     services.JobStatusPending,
			TotalItems: len(batchMedicos),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchMedicos)
		if err != nil {
			return fmt.Errorf("erro ao serializar batch: %v", err)
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
			return fmt.Errorf("erro ao enfileirar job: %v", err)
		}

		log.Printf("✓ Batch %d/%d enfileirado com %d médicos", i+1, totalBatches, len(batchMedicos))
	}

	return nil
}

func (h *WebSocketHandler) processUnifiedMunicipios(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	municipios, err := h.convertRowsToMunicipios(parsedData)
	if err != nil {
		return fmt.Errorf("erro ao converter dados para municípios: %v", err)
	}

	if len(municipios) == 0 {
		return fmt.Errorf("nenhum município válido encontrado nos dados")
	}

	totalItems := len(municipios)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processando %d municípios em %d batches para sessão %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchMunicipios := municipios[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "municipios",
			Status:     services.JobStatusPending,
			TotalItems: len(batchMunicipios),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchMunicipios)
		if err != nil {
			return fmt.Errorf("erro ao serializar batch: %v", err)
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
			return fmt.Errorf("erro ao enfileirar job: %v", err)
		}

		log.Printf("✓ Batch %d/%d enfileirado com %d municípios", i+1, totalBatches, len(batchMunicipios))
	}

	return nil
}

func (h *WebSocketHandler) processUnifiedEstados(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	estados, err := h.convertRowsToEstados(parsedData)
	if err != nil {
		return fmt.Errorf("erro ao converter dados para estados: %v", err)
	}

	if len(estados) == 0 {
		return fmt.Errorf("nenhum estado válido encontrado nos dados")
	}

	totalItems := len(estados)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processando %d estados em %d batches para sessão %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchEstados := estados[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "estados",
			Status:     services.JobStatusPending,
			TotalItems: len(batchEstados),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchEstados)
		if err != nil {
			return fmt.Errorf("erro ao serializar batch: %v", err)
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
			return fmt.Errorf("erro ao enfileirar job: %v", err)
		}

		log.Printf("✓ Batch %d/%d enfileirado com %d estados", i+1, totalBatches, len(batchEstados))
	}

	return nil
}

// processUnifiedCID10 processa dados CID10 usando dados unificados
func (h *WebSocketHandler) processUnifiedCID10(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	cid10s, err := h.convertRowsToCID10(parsedData)
	if err != nil {
		return fmt.Errorf("erro ao converter dados para CID10: %v", err)
	}

	if len(cid10s) == 0 {
		return fmt.Errorf("nenhum registro CID10 válido encontrado nos dados")
	}

	totalItems := len(cid10s)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processando %d registros CID10 em %d batches para sessão %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchCID10 := cid10s[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "cid10",
			Status:     services.JobStatusPending,
			TotalItems: len(batchCID10),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchCID10)
		if err != nil {
			return fmt.Errorf("erro ao serializar batch: %v", err)
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
			return fmt.Errorf("erro ao enfileirar job: %v", err)
		}

		log.Printf("✓ Batch %d/%d enfileirado com %d registros CID10", i+1, totalBatches, len(batchCID10))
	}

	return nil
}

// convertRowsToCID10 converte dados tabulares para CID10
func (h *WebSocketHandler) convertRowsToCID10(parsedData *parsers.ParsedData) ([]domain.Cid10, error) {
	var cid10s []domain.Cid10
	idCounter := 1

	// Se for arquivo com apenas uma coluna, assumir que cada linha é: CODIGO - DESCRIÇÃO
	if len(parsedData.Headers) == 1 {
		log.Printf("Processando arquivo CID10 com formato CODIGO - DESCRIÇÃO")
		
		for _, row := range parsedData.Rows {
			if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
				continue
			}

			linha := strings.TrimSpace(row[0])
			
			// Filtrar apenas linhas que começam com código válido (padrão: 1-3 letras + números)
			// Exemplos: C93, A01, B15, Z12, etc.
			if !h.isValidCIDCode(linha) {
				continue
			}

			// Separar código da descrição pelo primeiro "-"
			parts := strings.SplitN(linha, "-", 2)
			if len(parts) < 2 {
				continue
			}

			codigo := strings.TrimSpace(parts[0])
			descricao := strings.TrimSpace(parts[1])

			// Validar se código não está muito longo (códigos CID são curtos)
			if len(codigo) > 10 {
				continue
			}

			cid10 := domain.Cid10{
				ID:        idCounter,
				Codigo:    codigo,
				Descricao: descricao,
				Categoria: "",
				Grupo:     "",
			}

			cid10s = append(cid10s, cid10)
			idCounter++

			// Log primeiros 5 para debug
			if len(cid10s) <= 5 {
				log.Printf("✓ CID10 parseado: '%s' -> '%s'", codigo, descricao[:h.min(50, len(descricao))])
			}
		}
	} else {
		// Para arquivos com múltiplas colunas, usar lógica anterior
		return h.convertRowsToCID10Legacy(parsedData)
	}

	log.Printf("✓ Convertidos %d registros CID10 válidos dos dados", len(cid10s))
	
	return cid10s, nil
}

// isValidCIDCode verifica se uma linha começa com código CID válido
func (h *WebSocketHandler) isValidCIDCode(linha string) bool {
	if len(linha) < 3 {
		return false
	}
	
	trimmed := strings.TrimSpace(linha)
	if len(trimmed) < 3 {
		return false
	}

	// Primeiro caractere deve ser letra
	if !((trimmed[0] >= 'A' && trimmed[0] <= 'Z') || (trimmed[0] >= 'a' && trimmed[0] <= 'z')) {
		return false
	}

	// Procurar por números nos primeiros caracteres
	hasNumber := false
	for i := 1; i < h.min(6, len(trimmed)); i++ {
		if trimmed[i] >= '0' && trimmed[i] <= '9' {
			hasNumber = true
			break
		}
		// Se encontrar "-", parar busca
		if trimmed[i] == '-' {
			break
		}
	}

	return hasNumber
}

// min helper function
func (h *WebSocketHandler) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// convertRowsToCID10Legacy mantém a lógica anterior para arquivos com múltiplas colunas
func (h *WebSocketHandler) convertRowsToCID10Legacy(parsedData *parsers.ParsedData) ([]domain.Cid10, error) {
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var cid10s []domain.Cid10
	idCounter := 1

	for _, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}

		cid10 := domain.Cid10{}
		cid10.ID = idCounter

		// Código CID10
		codeFields := []string{"cid-10", "cid10", "codigo", "code", "classification"}
		for _, field := range codeFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				codigo := strings.TrimSpace(row[idx])
				if len(codigo) > 10 {
					continue
				}
				cid10.Codigo = codigo
				break
			}
		}

		// Descrição
		descFields := []string{"descricao", "description", "desc", "name", "nome", "titulo"}
		for _, field := range descFields {
			if idx, exists := headerMap[field]; exists && idx < len(row) {
				cid10.Descricao = strings.TrimSpace(row[idx])
				break
			}
		}

		// Só adicionar se tiver código válido
		if cid10.Codigo != "" && len(cid10.Codigo) <= 10 {
			cid10s = append(cid10s, cid10)
			idCounter++
		}
	}

	return cid10s, nil
}

func (h *WebSocketHandler) processFHIRData(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	// Detectar tipo de dados FHIR baseado nos headers
	dataType := h.detectFHIRDataType(parsedData)
	
	// Converter dados FHIR para objetos de domínio
	domainData, err := h.convertFHIRToDomain(parsedData, dataType)
	if err != nil {
		return fmt.Errorf("erro ao converter dados FHIR: %v", err)
	}
	
	// Processar baseado no tipo detectado
	switch dataType {
	case "patients", "pacientes":
		if patients, ok := domainData.([]domain.Paciente); ok {
			return h.enqueuePatientsJobs(sessionID, patients, batchSize)
		}
	case "practitioners", "medicos":
		if medicos, ok := domainData.([]domain.Medico); ok {
			// Processar médicos em batches
			for i := 0; i < len(medicos); i += batchSize {
				end := i + batchSize
				if end > len(medicos) {
					end = len(medicos)
				}
				batch := medicos[i:end]
				if err := h.submitMedicosBatch(sessionID, batch, i/batchSize+1); err != nil {
					return err
				}
			}
			return nil
		}
	case "organizations", "hospitais":
		if hospitais, ok := domainData.([]domain.Hospital); ok {
			// Processar hospitais em batches usando data service
			_, _, err := h.dataService.UpsertHospitais(hospitais)
			return err
		}
	case "locations", "municipios":
		if municipios, ok := domainData.([]domain.Municipio); ok {
			// Processar municípios em batches usando data service
			_, _, err := h.dataService.UpsertMunicipiosConcurrent(municipios)
			return err
		}
	default:
		return fmt.Errorf("tipo de dados FHIR não suportado: %s", dataType)
	}
	
	return fmt.Errorf("erro ao processar dados FHIR: tipo não reconhecido")
}

// detectFHIRDataType detecta o tipo de dados FHIR baseado nos headers
func (h *WebSocketHandler) detectFHIRDataType(parsedData *parsers.ParsedData) string {
	// Verificar se há resource_type nos headers
	for i, header := range parsedData.Headers {
		if strings.ToLower(header) == "resource_type" {
			// Verificar os valores na primeira linha
			if len(parsedData.Rows) > 0 && i < len(parsedData.Rows[0]) {
				resourceType := strings.ToLower(parsedData.Rows[0][i])
				switch resourceType {
				case "patient":
					return "patients"
				case "practitioner":
					return "practitioners"
				case "organization":
					return "organizations"
				case "location":
					return "locations"
				}
			}
		}
	}
	
	// Fallback: tentar detectar por campos específicos
	if h.hasFHIRField(parsedData, "patient_id") || h.hasFHIRField(parsedData, "cpf") {
		return "patients"
	}
	if h.hasFHIRField(parsedData, "practitioner_id") || h.hasFHIRField(parsedData, "especialidade") {
		return "practitioners"
	}
	if h.hasFHIRField(parsedData, "organization_id") || h.hasFHIRField(parsedData, "leitos_totais") {
		return "organizations"
	}
	if h.hasFHIRField(parsedData, "location_id") || h.hasFHIRField(parsedData, "codigo_ibge") {
		return "locations"
	}
	
	return "unknown"
}

// hasFHIRField verifica se um campo específico existe nos headers
func (h *WebSocketHandler) hasFHIRField(parsedData *parsers.ParsedData, fieldName string) bool {
	for _, header := range parsedData.Headers {
		if strings.ToLower(header) == strings.ToLower(fieldName) {
			return true
		}
	}
	return false
}

// convertFHIRToDomain converte dados FHIR parseados para objetos de domínio
func (h *WebSocketHandler) convertFHIRToDomain(parsedData *parsers.ParsedData, dataType string) (interface{}, error) {
	// Criar um parser unificado para conversão
	unifiedParser := parsers.NewUnifiedParser()
	
	// Converter usando o parser unificado
	return unifiedParser.ConvertToDataType(parsedData, dataType)
}

// processHL7Data processa dados HL7
func (h *WebSocketHandler) processHL7Data(sessionID string, parsedData *parsers.ParsedData, batchSize int) error {
	// Detectar tipo de dados HL7 baseado nos headers
	dataType := h.detectHL7DataType(parsedData)
	
	// Converter dados HL7 para objetos de domínio
	domainData, err := h.convertHL7ToDomain(parsedData, dataType)
	if err != nil {
		return fmt.Errorf("erro ao converter dados HL7: %v", err)
	}
	
	// Processar baseado no tipo detectado
	switch dataType {
	case "patients", "pacientes":
		if patients, ok := domainData.([]domain.Paciente); ok {
			return h.enqueuePatientsJobs(sessionID, patients, batchSize)
		}
	case "practitioners", "medicos":
		if medicos, ok := domainData.([]domain.Medico); ok {
			// Processar médicos em batches
			for i := 0; i < len(medicos); i += batchSize {
				end := i + batchSize
				if end > len(medicos) {
					end = len(medicos)
				}
				batch := medicos[i:end]
				if err := h.submitMedicosBatch(sessionID, batch, i/batchSize+1); err != nil {
					return err
				}
			}
			return nil
		}
	case "organizations", "hospitais":
		if hospitais, ok := domainData.([]domain.Hospital); ok {
			// Processar hospitais em batches usando data service
			_, _, err := h.dataService.UpsertHospitais(hospitais)
			return err
		}
	default:
		return fmt.Errorf("tipo de dados HL7 não suportado: %s", dataType)
	}
	
	return fmt.Errorf("erro ao processar dados HL7: tipo não reconhecido")
}

// detectHL7DataType detecta o tipo de dados HL7 baseado nos headers
func (h *WebSocketHandler) detectHL7DataType(parsedData *parsers.ParsedData) string {
	// Verificar se há message_type nos headers
	for i, header := range parsedData.Headers {
		if strings.ToLower(header) == "message_type" {
			// Verificar os valores na primeira linha
			if len(parsedData.Rows) > 0 && i < len(parsedData.Rows[0]) {
				messageType := strings.ToLower(parsedData.Rows[0][i])
				// Detectar tipo baseado no prefixo da mensagem
				if strings.HasPrefix(messageType, "adt^") {
					return "patients"
				}
				if strings.HasPrefix(messageType, "mfn^") {
					// Para MFN, verificar se é practitioner ou organization
					// Por simplicidade, assumir practitioners por padrão
					return "practitioners"
				}
				if strings.HasPrefix(messageType, "org^") {
					return "organizations"
				}
			}
		}
	}
	
	// Fallback: tentar detectar por campos específicos
	if h.hasHL7Field(parsedData, "patient_id") || h.hasHL7Field(parsedData, "patient_name") {
		return "patients"
	}
	if h.hasHL7Field(parsedData, "practitioner_id") || h.hasHL7Field(parsedData, "especialidade") {
		return "practitioners"
	}
	if h.hasHL7Field(parsedData, "organization_id") || h.hasHL7Field(parsedData, "leitos_totais") {
		return "organizations"
	}
	
	return "patients" // Default para pacientes
}

// hasHL7Field verifica se um campo específico existe nos headers HL7
func (h *WebSocketHandler) hasHL7Field(parsedData *parsers.ParsedData, fieldName string) bool {
	for _, header := range parsedData.Headers {
		if strings.ToLower(header) == strings.ToLower(fieldName) {
			return true
		}
	}
	return false
}

// convertHL7ToDomain converte dados HL7 parseados para objetos de domínio
func (h *WebSocketHandler) convertHL7ToDomain(parsedData *parsers.ParsedData, dataType string) (interface{}, error) {
	// Criar um parser unificado para conversão
	unifiedParser := parsers.NewUnifiedParser()
	
	// Converter usando o parser unificado
	return unifiedParser.ConvertToDataType(parsedData, dataType)
}

func (h *WebSocketHandler) processUploadedDataLegacy(sessionID, fileType, csvData string) error {
	// Configuração de batch por tipo
	batchSizes := map[string]int{
		"estados":    getEnvAsInt("BATCH_SIZE_ESTADOS", 50),
		"municipios": getEnvAsInt("BATCH_SIZE_MUNICIPIOS", 200),
		"medicos":    getEnvAsInt("BATCH_SIZE_MEDICOS", 300),
		"hospitais":  getEnvAsInt("BATCH_SIZE_HOSPITAIS", 150),
		"pacientes":  getEnvAsInt("BATCH_SIZE_PACIENTES", 250),
		"cid10":      getEnvAsInt("BATCH_SIZE_CID10", 400),
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
			ID:         jobID,
			Type:       "municipios",
			Status:     services.JobStatusPending,
			TotalItems: len(batchMunicipios),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
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
	return h.processMedicosStreaming(sessionID, csvData, batchSize)
}

// processMedicosStreaming processa CSV em streaming, criando batches conforme lê os dados
func (h *WebSocketHandler) processMedicosStreaming(sessionID, csvData string, batchSize int) error {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Ler cabeçalho
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	// Criar mapa de índices dos cabeçalhos
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	var currentBatch []domain.Medico
	batchNumber := 1
	totalProcessed := 0

	log.Printf("Starting streaming processing for session %s with batch size %d", sessionID, batchSize)

	// Processar linha por linha
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			continue
		}

		// Parse do médico
		medico := domain.Medico{}

		// UUID do médico (da coluna codigo)
		if idx, exists := headerMap["codigo"]; exists && idx < len(record) {
			uuidStr := strings.TrimSpace(record[idx])
			medico.UUID, _ = uuid.Parse(uuidStr)
		}

		// Nome (múltiplas variações)
		nomeColumns := []string{"nome_completo", "nome", "name", "medico", "doctor", "nome_medico", "doctor_name"}
		for _, col := range nomeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				medico.Nome = strings.TrimSpace(record[idx])
				break
			}
		}

		// Especialidade
		especialidadeColumns := []string{"especialidade", "specialty", "especialization", "area", "speciality"}
		for _, col := range especialidadeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				medico.Especialidade = strings.TrimSpace(record[idx])
				break
			}
		}

		// Código do município
		municipioColumns := []string{"cidade", "cod_municipio", "codigo_municipio", "municipio_id", "city", "city_code", "municipality_code"}
		for _, col := range municipioColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				codMunicipio := strings.TrimSpace(record[idx])
				if codMunicipio != "" {
					medico.CodMunicipio = codMunicipio
					break
				}
			}
		}

		// Só adicionar se tiver pelo menos nome
		if medico.Nome != "" {
			currentBatch = append(currentBatch, medico)
			totalProcessed++

			// Se atingiu o tamanho do batch, processar
			if len(currentBatch) >= batchSize {
				err := h.submitMedicosBatch(sessionID, currentBatch, batchNumber)
				if err != nil {
					log.Printf("Error submitting batch %d: %v", batchNumber, err)
				} else {
					log.Printf("✓ Batch %d submitted with %d doctors", batchNumber, len(currentBatch))
				}

				// Reset batch
				currentBatch = []domain.Medico{}
				batchNumber++
			}
		}
	}

	// Processar último batch se houver registros restantes
	if len(currentBatch) > 0 {
		err := h.submitMedicosBatch(sessionID, currentBatch, batchNumber)
		if err != nil {
			log.Printf("Error submitting final batch %d: %v", batchNumber, err)
		} else {
			log.Printf("✓ Final batch %d submitted with %d doctors", batchNumber, len(currentBatch))
		}
	}

	log.Printf("Streaming processing completed. Total doctors processed: %d in %d batches", totalProcessed, batchNumber)
	return nil
}

// submitMedicosBatch submete um batch de médicos para a fila Redis
func (h *WebSocketHandler) submitMedicosBatch(sessionID string, medicos []domain.Medico, batchNumber int) error {
	jobID := uuid.New().String()
	job := &services.Job{
		ID:         jobID,
		Type:       "medicos",
		Status:     services.JobStatusPending,
		TotalItems: len(medicos),
		CreatedAt:  time.Now(),
		SessionID:  sessionID,
	}

	batchDataBytes, err := json.Marshal(medicos)
	if err != nil {
		return fmt.Errorf("failed to marshal batch data: %v", err)
	}

	batchData := services.BatchData{
		BatchNumber:  batchNumber,
		TotalBatches: -1, // Não sabemos o total antecipadamente no streaming
		Data:         batchDataBytes,
	}

	data, _ := json.Marshal(batchData)
	job.Data = data

	err = h.redisService.EnqueueJob(job)
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %v", err)
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
			ID:         jobID,
			Type:       "estados",
			Status:     services.JobStatusPending,
			TotalItems: len(batchEstados),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
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

// processHospitais processa dados de hospitais
func (h *WebSocketHandler) processHospitais(sessionID, csvData string, batchSize int) error {
	hospitais, err := h.parseCSVHospitaisFromString(csvData)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(hospitais) == 0 {
		return fmt.Errorf("no valid hospitals found in CSV")
	}

	totalItems := len(hospitais)
	totalBatches := (totalItems + batchSize - 1) / batchSize

	log.Printf("Processing %d hospitals in %d batches for session %s", totalItems, totalBatches, sessionID)

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		batchHospitais := hospitais[start:end]

		jobID := uuid.New().String()
		job := &services.Job{
			ID:         jobID,
			Type:       "hospitais",
			Status:     services.JobStatusPending,
			TotalItems: len(batchHospitais),
			CreatedAt:  time.Now(),
			SessionID:  sessionID,
		}

		batchDataBytes, err := json.Marshal(batchHospitais)
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

		log.Printf("Enqueued batch %d/%d with %d hospitals", i+1, totalBatches, len(batchHospitais))
	}

	return nil
}

// parseCSVHospitaisFromString parses CSV data from string
func (h *WebSocketHandler) parseCSVHospitaisFromString(csvData string) ([]domain.Hospital, error) {
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

	log.Printf("CSV headers para hospitais: %v", headers[0])

	headerMap := make(map[string]int)
	for i, header := range headers[0] {
		cleanHeader := strings.ToLower(strings.TrimSpace(header))
		headerMap[cleanHeader] = i
	}

	var hospitais []domain.Hospital

	for i := 1; i < len(headers); i++ {
		record := headers[i]
		hospital := domain.Hospital{}

		// UUID do hospital (da coluna codigo)
		if idx, exists := headerMap["codigo"]; exists && idx < len(record) {
			uuidStr := strings.TrimSpace(record[idx])
			if uuidStr != "" {
				hospital.UUID, _ = uuid.Parse(uuidStr)
			} else {
				hospital.UUID = uuid.New()
			}
		} else {
			hospital.UUID = uuid.New()
		}

		// Nome do hospital
		nomeColumns := []string{"nome", "hospital_nome", "hospital", "name", "hospital_name"}
		for _, col := range nomeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				hospital.Nome = strings.TrimSpace(record[idx])
				break
			}
		}

		// CEP
		cepColumns := []string{"cep", "zipcode", "postal_code"}
		for _, col := range cepColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				hospital.CEP = strings.TrimSpace(record[idx])
				break
			}
		}

		// Especialidades (mantém como string com separador ;)
		especialidadeColumns := []string{"especialidades", "specialties", "services", "servicos"}
		for _, col := range especialidadeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				hospital.Especialidades = strings.TrimSpace(record[idx])
				break
			}
		}

		// Leitos totais
		leitosColumns := []string{"leitos_totais", "leitos", "beds", "total_beds"}
		for _, col := range leitosColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				if leitos, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
					hospital.LeitosTotais = leitos
				}
				break
			}
		}

		// Código do município/cidade
		municipioColumns := []string{"cidade", "cod_municipio", "codigo_municipio", "municipio_id", "city", "city_code", "municipality_code"}
		for _, col := range municipioColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				hospital.CodMunicipio = strings.TrimSpace(record[idx])
				break
			}
		}

		// Bairro
		bairroColumns := []string{"bairro", "district", "neighborhood"}
		for _, col := range bairroColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				hospital.Bairro = strings.TrimSpace(record[idx])
				break
			}
		}

		// Log detalhado do parsing
		log.Printf("Registro %d - Nome: '%s', CEP: '%s', Especialidades: '%s', Leitos: %d", 
			i, hospital.Nome, hospital.CEP, hospital.Especialidades, hospital.LeitosTotais)

		// Só adicionar se tiver pelo menos nome
		if hospital.Nome != "" {
			hospitais = append(hospitais, hospital)
			log.Printf("✓ Hospital %d adicionado: %s", i, hospital.Nome)
		} else {
			log.Printf("✗ Hospital %d rejeitado - nome vazio", i)
		}
	}

	log.Printf("Total de hospitais parseados: %d", len(hospitais))
	return hospitais, nil
}

func (h *WebSocketHandler) processPacientes(sessionID, csvData string, batchSize int) error {
	return h.processPacientesStreaming(sessionID, csvData, batchSize)
}

// processPacientesStreaming processa CSV em streaming, criando batches conforme lê os dados
func (h *WebSocketHandler) processPacientesStreaming(sessionID, csvData string, batchSize int) error {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Ler cabeçalho
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	// Criar mapa de índices dos cabeçalhos
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	log.Printf("🔍 DEBUG Pacientes - Headers: %v", headers)
	log.Printf("🔍 DEBUG Pacientes - Header map: %v", headerMap)

	var currentBatch []domain.Paciente
	batchNumber := 1
	totalProcessed := 0

	log.Printf("Starting streaming processing for pacientes session %s with batch size %d", sessionID, batchSize)

	// Processar linha por linha
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			continue
		}

		// Parse do paciente
		paciente := domain.Paciente{}

		// ID do paciente (UUID)
		if idx, exists := headerMap["id"]; exists && idx < len(record) {
			uuidStr := strings.TrimSpace(record[idx])
			if uuidStr != "" {
				if parsedUUID, err := uuid.Parse(uuidStr); err == nil {
					paciente.ID = parsedUUID
				} else {
					paciente.ID = uuid.New()
				}
			} else {
				paciente.ID = uuid.New()
			}
		} else {
			paciente.ID = uuid.New()
		}

		// CPF
		cpfColumns := []string{"cpf", "documento", "doc"}
		for _, col := range cpfColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				paciente.CPF = strings.TrimSpace(record[idx])
				break
			}
		}

		// Nome
		nomeColumns := []string{"nome", "nome_completo", "name", "full_name", "patient_name"}
		for _, col := range nomeColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				paciente.Nome = strings.TrimSpace(record[idx])
				break
			}
		}

		// Gênero
		generoColumns := []string{"genero", "gender", "sex", "sexo"}
		for _, col := range generoColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				genero := strings.TrimSpace(record[idx])
				// Normalizar para M/F
				if len(genero) > 0 {
					switch strings.ToUpper(genero[:1]) {
					case "M":
						paciente.Genero = "M"
					case "F":
						paciente.Genero = "F"
					default:
						paciente.Genero = genero[:1] // Manter o primeiro caractere
					}
				}
				break
			}
		}

		// Código do município
		municipioColumns := []string{"cod_municipio", "codigo_municipio", "municipio_id", "city_code", "cidade"}
		for _, col := range municipioColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				paciente.CodMunicipio = strings.TrimSpace(record[idx])
				break
			}
		}

		// Bairro
		bairroColumns := []string{"bairro", "district", "neighborhood", "endereco"}
		for _, col := range bairroColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				paciente.Bairro = strings.TrimSpace(record[idx])
				break
			}
		}

		// Convênio
		convenioColumns := []string{"convenio", "insurance", "plano_saude", "plano"}
		for _, col := range convenioColumns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				paciente.Convenio = strings.TrimSpace(record[idx])
				break
			}
		}

		// CID10
		cid10Columns := []string{"cid10", "cid-10", "cid_10", "diagnosis", "diagnostico"}
		for _, col := range cid10Columns {
			if idx, exists := headerMap[col]; exists && idx < len(record) {
				paciente.CID10 = strings.TrimSpace(record[idx])
				break
			}
		}

		// Só adicionar se tiver pelo menos CPF e nome
		if paciente.CPF != "" && paciente.Nome != "" {
			currentBatch = append(currentBatch, paciente)
			totalProcessed++

			log.Printf("✅ DEBUG Pacientes - Adicionado: CPF=%s, Nome=%s, Genero=%s",
				paciente.CPF, paciente.Nome, paciente.Genero)

			// Se atingiu o tamanho do batch, processar
			if len(currentBatch) >= batchSize {
				err := h.submitPacientesBatch(sessionID, currentBatch, batchNumber)
				if err != nil {
					log.Printf("Error submitting batch %d: %v", batchNumber, err)
				} else {
					log.Printf("✓ Batch %d submitted with %d pacientes", batchNumber, len(currentBatch))
				}

				// Reset batch
				currentBatch = []domain.Paciente{}
				batchNumber++
			}
		} else {
			log.Printf("⚠️ DEBUG Pacientes - Linha ignorada: CPF='%s', Nome='%s'", paciente.CPF, paciente.Nome)
		}
	}

	// Processar último batch se houver registros restantes
	if len(currentBatch) > 0 {
		err := h.submitPacientesBatch(sessionID, currentBatch, batchNumber)
		if err != nil {
			log.Printf("Error submitting final batch %d: %v", batchNumber, err)
		} else {
			log.Printf("✓ Final batch %d submitted with %d pacientes", batchNumber, len(currentBatch))
		}
	}

	log.Printf("Streaming processing completed. Total pacientes processed: %d in %d batches", totalProcessed, batchNumber)
	return nil
}

// submitPacientesBatch submete um batch de pacientes para a fila Redis
func (h *WebSocketHandler) submitPacientesBatch(sessionID string, pacientes []domain.Paciente, batchNumber int) error {
	jobID := uuid.New().String()
	job := &services.Job{
		ID:         jobID,
		Type:       "pacientes",
		Status:     services.JobStatusPending,
		TotalItems: len(pacientes),
		CreatedAt:  time.Now(),
		SessionID:  sessionID,
	}

	batchDataBytes, err := json.Marshal(pacientes)
	if err != nil {
		return fmt.Errorf("failed to marshal batch data: %v", err)
	}

	batchData := services.BatchData{
		BatchNumber:  batchNumber,
		TotalBatches: -1, // Não sabemos o total antecipadamente no streaming
		Data:         batchDataBytes,
	}

	data, _ := json.Marshal(batchData)
	job.Data = data

	err = h.redisService.EnqueueJob(job)
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %v", err)
	}

	return nil
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
