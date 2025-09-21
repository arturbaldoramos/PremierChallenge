package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"example.com/m/v2/internal/domain"
	"example.com/m/v2/internal/parsers"
	"example.com/m/v2/internal/services"
	"gorm.io/gorm"
)

type UploadHandler struct {
	dataService    *services.DataService
	csvParser      *parsers.CSVParser
	unifiedParser  *parsers.UnifiedParser
}

type UploadResponse struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message"`
	FileType        string   `json:"file_type"`
	FileName        string   `json:"file_name"`
	TotalRecords    int      `json:"total_records"`
	InsertedRecords int      `json:"inserted_records"`
	UpdatedRecords  int      `json:"updated_records"`
	Errors          []string `json:"errors,omitempty"`
}

func NewUploadHandler(db *gorm.DB) *UploadHandler {
	return &UploadHandler{
		dataService:   services.NewDataService(db),
		csvParser:     parsers.NewCSVParser(),
		unifiedParser: parsers.NewUnifiedParser(),
	}
}

func (h *UploadHandler) UploadEstados(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB limit
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get file", err)
		return
	}
	defer file.Close()

	// Detect file type
	fileInfo, err := parsers.DetectFileType(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to detect file type", err)
		return
	}

	// Ler arquivo para bytes
	fileBytes, err := h.readFileToBytes(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read file", err)
		return
	}

	// Usar UnifiedParser para processar qualquer formato
	parsedData, format, err := h.unifiedParser.ParseFile(fileBytes, "estados")
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse file", err)
		return
	}

	// Converter para objetos de domínio
	domainData, err := h.unifiedParser.ConvertToDataType(parsedData, "estados")
	if err != nil {
		h.sendErrorResponse(w, "Failed to convert to Estados objects", err)
		return
	}

	estados, ok := domainData.([]domain.Estado)
	if !ok {
		h.sendErrorResponse(w, "Failed to convert to Estados objects", fmt.Errorf("invalid data type"))
		return
	}

	// Salvar no banco
	inserted, updated, err := h.dataService.UpsertEstados(estados)
	if err != nil {
		h.sendErrorResponse(w, "Failed to save estados", err)
		return
	}

	response := UploadResponse{
		Success:         true,
		Message:         "Estados uploaded successfully",
		FileType:        format.String(),
		FileName:        fileInfo.Name,
		TotalRecords:    len(estados),
		InsertedRecords: inserted,
		UpdatedRecords:  updated,
	}

	h.sendSuccessResponse(w, response)
}

func (h *UploadHandler) UploadMunicipios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get file", err)
		return
	}
	defer file.Close()

	fileInfo, err := parsers.DetectFileType(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to detect file type", err)
		return
	}

	// Ler arquivo para bytes
	fileBytes, err := h.readFileToBytes(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read file", err)
		return
	}

	// Usar UnifiedParser para processar qualquer formato
	parsedData, format, err := h.unifiedParser.ParseFile(fileBytes, "municipios")
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse file", err)
		return
	}

	// Converter para objetos de domínio
	domainData, err := h.unifiedParser.ConvertToDataType(parsedData, "municipios")
	if err != nil {
		h.sendErrorResponse(w, "Failed to convert to Municipios objects", err)
		return
	}

	municipios, ok := domainData.([]domain.Municipio)
	if !ok {
		h.sendErrorResponse(w, "Failed to convert to Municipios objects", fmt.Errorf("invalid data type"))
		return
	}

	// Usar versão concorrente para melhor performance
	inserted, updated, err := h.dataService.UpsertMunicipiosConcurrent(municipios)
	if err != nil {
		h.sendErrorResponse(w, "Failed to save municipios", err)
		return
	}

	response := UploadResponse{
		Success:         true,
		Message:         "Municipios uploaded successfully",
		FileType:        format.String(),
		FileName:        fileInfo.Name,
		TotalRecords:    len(municipios),
		InsertedRecords: inserted,
		UpdatedRecords:  updated,
	}

	h.sendSuccessResponse(w, response)
}

func (h *UploadHandler) UploadHospitais(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get file", err)
		return
	}
	defer file.Close()

	fileInfo, err := parsers.DetectFileType(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to detect file type", err)
		return
	}

	// Ler arquivo para bytes
	fileBytes, err := h.readFileToBytes(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read file", err)
		return
	}

	// Usar UnifiedParser para processar qualquer formato
	parsedData, format, err := h.unifiedParser.ParseFile(fileBytes, "hospitais")
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse file", err)
		return
	}

	// Converter para objetos de domínio
	domainData, err := h.unifiedParser.ConvertToDataType(parsedData, "hospitais")
	if err != nil {
		h.sendErrorResponse(w, "Failed to convert to Hospitais objects", err)
		return
	}

	hospitais, ok := domainData.([]domain.Hospital)
	if !ok {
		h.sendErrorResponse(w, "Failed to convert to Hospitais objects", fmt.Errorf("invalid data type"))
		return
	}

	inserted, updated, err := h.dataService.UpsertHospitais(hospitais)
	if err != nil {
		h.sendErrorResponse(w, "Failed to save hospitais", err)
		return
	}

	response := UploadResponse{
		Success:         true,
		Message:         "Hospitais uploaded successfully",
		FileType:        format.String(),
		FileName:        fileInfo.Name,
		TotalRecords:    len(hospitais),
		InsertedRecords: inserted,
		UpdatedRecords:  updated,
	}

	h.sendSuccessResponse(w, response)
}

func (h *UploadHandler) UploadPacientes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get file", err)
		return
	}
	defer file.Close()

	fileInfo, err := parsers.DetectFileType(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to detect file type", err)
		return
	}

	// Ler arquivo para bytes
	fileBytes, err := h.readFileToBytes(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read file", err)
		return
	}

	// Usar UnifiedParser para processar qualquer formato
	parsedData, format, err := h.unifiedParser.ParseFile(fileBytes, "pacientes")
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse file", err)
		return
	}

	// Converter para objetos de domínio
	domainData, err := h.unifiedParser.ConvertToDataType(parsedData, "pacientes")
	if err != nil {
		h.sendErrorResponse(w, "Failed to convert to Pacientes objects", err)
		return
	}

	pacientes, ok := domainData.([]domain.Paciente)
	if !ok {
		h.sendErrorResponse(w, "Failed to convert to Pacientes objects", fmt.Errorf("invalid data type"))
		return
	}

	inserted, updated, err := h.dataService.UpsertPacientes(pacientes)
	if err != nil {
		h.sendErrorResponse(w, "Failed to save pacientes", err)
		return
	}

	response := UploadResponse{
		Success:         true,
		Message:         "Pacientes uploaded successfully",
		FileType:        format.String(),
		FileName:        fileInfo.Name,
		TotalRecords:    len(pacientes),
		InsertedRecords: inserted,
		UpdatedRecords:  updated,
	}

	h.sendSuccessResponse(w, response)
}

func (h *UploadHandler) UploadMedicos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get file", err)
		return
	}
	defer file.Close()

	fileInfo, err := parsers.DetectFileType(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to detect file type", err)
		return
	}

	// Ler arquivo para bytes
	fileBytes, err := h.readFileToBytes(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read file", err)
		return
	}

	// Usar UnifiedParser para processar qualquer formato
	parsedData, format, err := h.unifiedParser.ParseFile(fileBytes, "medicos")
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse file", err)
		return
	}

	// Converter para objetos de domínio
	domainData, err := h.unifiedParser.ConvertToDataType(parsedData, "medicos")
	if err != nil {
		h.sendErrorResponse(w, "Failed to convert to Medicos objects", err)
		return
	}

	medicos, ok := domainData.([]domain.Medico)
	if !ok {
		h.sendErrorResponse(w, "Failed to convert to Medicos objects", fmt.Errorf("invalid data type"))
		return
	}

	// Usar versão concorrente para melhor performance
	inserted, updated, err := h.dataService.UpsertMedicosConcurrent(medicos)
	if err != nil {
		h.sendErrorResponse(w, "Failed to save medicos", err)
		return
	}

	response := UploadResponse{
		Success:         true,
		Message:         "Medicos uploaded successfully",
		FileType:        format.String(),
		FileName:        fileInfo.Name,
		TotalRecords:    len(medicos),
		InsertedRecords: inserted,
		UpdatedRecords:  updated,
	}

	h.sendSuccessResponse(w, response)
}

func (h *UploadHandler) UploadCID10(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get file", err)
		return
	}
	defer file.Close()

	fileInfo, err := parsers.DetectFileType(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to detect file type", err)
		return
	}

	// Ler arquivo para bytes
	fileBytes, err := h.readFileToBytes(fileHeader)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read file", err)
		return
	}

	// Usar UnifiedParser para processar qualquer formato
	parsedData, format, err := h.unifiedParser.ParseFile(fileBytes, "cid10")
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse file", err)
		return
	}

	// Converter para objetos de domínio
	domainData, err := h.unifiedParser.ConvertToDataType(parsedData, "cid10")
	if err != nil {
		h.sendErrorResponse(w, "Failed to convert to CID10 objects", err)
		return
	}

	cid10s, ok := domainData.([]domain.Cid10)
	if !ok {
		h.sendErrorResponse(w, "Failed to convert to CID10 objects", fmt.Errorf("invalid data type"))
		return
	}

	// Salvar no banco
	inserted, updated, err := h.dataService.UpsertCID10(cid10s)
	if err != nil {
		h.sendErrorResponse(w, "Failed to save CID10", err)
		return
	}

	response := UploadResponse{
		Success:         true,
		Message:         "CID10 uploaded successfully",
		FileType:        format.String(),
		FileName:        fileInfo.Name,
		TotalRecords:    len(cid10s),
		InsertedRecords: inserted,
		UpdatedRecords:  updated,
	}

	h.sendSuccessResponse(w, response)
}

// readFileToBytes lê um arquivo multipart para bytes
func (h *UploadHandler) readFileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Ler todo o arquivo para bytes
	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *UploadHandler) sendSuccessResponse(w http.ResponseWriter, response UploadResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UploadHandler) sendErrorResponse(w http.ResponseWriter, message string, err error) {
	response := UploadResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Errors = []string{err.Error()}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}
