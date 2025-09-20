package parsers

import (
	"bytes"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type FileType string

const (
	CSV  FileType = "csv"
	XLSX FileType = "xlsx"
	XML  FileType = "xml"
)

type FileInfo struct {
	Type     FileType
	Name     string
	Size     int64
	MimeType string
}

// DetectFileType detects the file type based on extension and content
func DetectFileType(fileHeader *multipart.FileHeader) (FileInfo, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()

	// Read first few bytes for content detection
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	content := buffer[:n]

	// Get file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	var fileType FileType

	// Detect based on extension and magic bytes
	switch ext {
	case ".csv":
		fileType = CSV
	case ".xlsx", ".xls":
		if bytes.Contains(content, []byte("PK")) || bytes.Contains(content, []byte("\xd0\xcf\x11\xe0")) {
			fileType = XLSX
		} else {
			fileType = CSV // fallback for malformed xlsx
		}
	case ".xml":
		if bytes.Contains(content, []byte("<?xml")) || bytes.Contains(content, []byte("<")) {
			fileType = XML
		} else {
			fileType = CSV // fallback
		}
	default:
		// Try to detect by content
		if bytes.Contains(content, []byte("<?xml")) || bytes.Contains(content, []byte("<")) {
			fileType = XML
		} else if bytes.Contains(content, []byte("PK")) {
			fileType = XLSX
		} else {
			fileType = CSV // default fallback
		}
	}

	return FileInfo{
		Type:     fileType,
		Name:     fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: fileHeader.Header.Get("Content-Type"),
	}, nil
}