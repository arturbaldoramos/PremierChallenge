package parsers

import (
	"mime/multipart"

	"example.com/m/v2/internal/domain"
)

// Parser interface for different file types and domains
type Parser interface {
	ParseEstados(file *multipart.FileHeader) ([]domain.Estado, error)
	ParseMunicipios(file *multipart.FileHeader) ([]domain.Municipio, error)
	ParseHospitais(file *multipart.FileHeader) ([]domain.Hospital, error)
	ParsePacientes(file *multipart.FileHeader) ([]domain.Paciente, error)
	ParseMedicos(file *multipart.FileHeader) ([]domain.Medico, error)
	ParseCID10(file *multipart.FileHeader) ([]domain.Cid10, error)
}

// ParseResult represents the result of a parsing operation
type ParseResult struct {
	TotalRecords    int
	ParsedRecords   int
	SkippedRecords  int
	DuplicateRecords int
	Errors          []ParseError
}

// ParseError represents an error during parsing
type ParseError struct {
	Row     int
	Column  string
	Message string
	Data    interface{}
}

// ChunkProcessor interface for handling large files
type ChunkProcessor interface {
	ProcessChunk(data []byte, chunkIndex int) error
	GetChunkSize() int
	IsCompleted() bool
}