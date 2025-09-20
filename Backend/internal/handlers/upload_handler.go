package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/m/v2/internal/parsers"
	"example.com/m/v2/internal/services"
	"gorm.io/gorm"
)

type UploadHandler struct {
	dataService *services.DataService
	csvParser   *parsers.CSVParser
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
		dataService: services.NewDataService(db),
		csvParser:   parsers.NewCSVParser(),
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

	// Parse based on file type
	var estados []parsers.EstadoCSV
	switch fileInfo.Type {
	case parsers.CSV:
		domainEstados, err := h.csvParser.ParseEstados(fileHeader)
		if err != nil {
			h.sendErrorResponse(w, "Failed to parse CSV file", err)
			return
		}

		// Convert to domain objects
		for _, estado := range domainEstados {
			estados = append(estados, parsers.EstadoCSV{
				Codigo:            estado.Codigo,
				UnidadeFederativa: estado.UnidadeFederativa,
				Nome:              estado.Nome,
				Regiao:            estado.Regiao,
				Latitude:          estado.Latitude,
				Longitude:         estado.Longitude,
			})
		}

		// Save to database with duplicate checking
		inserted, updated, err := h.dataService.UpsertEstados(domainEstados)
		if err != nil {
			h.sendErrorResponse(w, "Failed to save estados", err)
			return
		}

		response := UploadResponse{
			Success:         true,
			Message:         "Estados uploaded successfully",
			FileType:        string(fileInfo.Type),
			FileName:        fileInfo.Name,
			TotalRecords:    len(domainEstados),
			InsertedRecords: inserted,
			UpdatedRecords:  updated,
		}

		h.sendSuccessResponse(w, response)

	default:
		h.sendErrorResponse(w, fmt.Sprintf("File type %s not supported yet", fileInfo.Type), nil)
	}
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

	switch fileInfo.Type {
	case parsers.CSV:
		// Usar versão streaming para economizar memória
		municipios, err := h.csvParser.ParseMunicipiosStreaming(fileHeader, 1000)
		if err != nil {
			h.sendErrorResponse(w, "Failed to parse CSV file", err)
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
			FileType:        string(fileInfo.Type),
			FileName:        fileInfo.Name,
			TotalRecords:    len(municipios),
			InsertedRecords: inserted,
			UpdatedRecords:  updated,
		}

		h.sendSuccessResponse(w, response)

	default:
		h.sendErrorResponse(w, fmt.Sprintf("File type %s not supported yet", fileInfo.Type), nil)
	}
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

	switch fileInfo.Type {
	case parsers.CSV:
		hospitais, err := h.csvParser.ParseHospitais(fileHeader)
		if err != nil {
			h.sendErrorResponse(w, "Failed to parse CSV file", err)
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
			FileType:        string(fileInfo.Type),
			FileName:        fileInfo.Name,
			TotalRecords:    len(hospitais),
			InsertedRecords: inserted,
			UpdatedRecords:  updated,
		}

		h.sendSuccessResponse(w, response)

	default:
		h.sendErrorResponse(w, fmt.Sprintf("File type %s not supported yet", fileInfo.Type), nil)
	}
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

	switch fileInfo.Type {
	case parsers.CSV:
		pacientes, err := h.csvParser.ParsePacientes(fileHeader)
		if err != nil {
			h.sendErrorResponse(w, "Failed to parse CSV file", err)
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
			FileType:        string(fileInfo.Type),
			FileName:        fileInfo.Name,
			TotalRecords:    len(pacientes),
			InsertedRecords: inserted,
			UpdatedRecords:  updated,
		}

		h.sendSuccessResponse(w, response)

	default:
		h.sendErrorResponse(w, fmt.Sprintf("File type %s not supported yet", fileInfo.Type), nil)
	}
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

	switch fileInfo.Type {
	case parsers.CSV:
		// Usar versão streaming para economizar memória
		medicos, err := h.csvParser.ParseMedicosStreaming(fileHeader, 1000)
		if err != nil {
			h.sendErrorResponse(w, "Failed to parse CSV file", err)
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
			FileType:        string(fileInfo.Type),
			FileName:        fileInfo.Name,
			TotalRecords:    len(medicos),
			InsertedRecords: inserted,
			UpdatedRecords:  updated,
		}

		h.sendSuccessResponse(w, response)

	default:
		h.sendErrorResponse(w, fmt.Sprintf("File type %s not supported yet", fileInfo.Type), nil)
	}
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

	switch fileInfo.Type {
	case parsers.CSV:
		cid10s, err := h.csvParser.ParseCID10(fileHeader)
		if err != nil {
			h.sendErrorResponse(w, "Failed to parse CSV file", err)
			return
		}

		inserted, updated, err := h.dataService.UpsertCID10(cid10s)
		if err != nil {
			h.sendErrorResponse(w, "Failed to save CID10", err)
			return
		}

		response := UploadResponse{
			Success:         true,
			Message:         "CID10 uploaded successfully",
			FileType:        string(fileInfo.Type),
			FileName:        fileInfo.Name,
			TotalRecords:    len(cid10s),
			InsertedRecords: inserted,
			UpdatedRecords:  updated,
		}

		h.sendSuccessResponse(w, response)

	default:
		h.sendErrorResponse(w, fmt.Sprintf("File type %s not supported yet", fileInfo.Type), nil)
	}
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
