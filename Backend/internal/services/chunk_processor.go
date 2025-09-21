package services

import (
	"bufio"
	"bytes"
	"io"
	"mime/multipart"
)

const (
	DefaultChunkSize = 1024 * 1024 // 1MB chunks
	MaxChunkSize     = 10 * 1024 * 1024 // 10MB max chunk size
)

type ChunkProcessor struct {
	chunkSize     int
	totalSize     int64
	processedSize int64
	completed     bool
}

type ChunkData struct {
	Index   int
	Data    []byte
	IsLast  bool
	Offset  int64
	Size    int
}

type ChunkResult struct {
	Index        int
	RecordsFound int
	Errors       []error
	Success      bool
}

func NewChunkProcessor(chunkSize int) *ChunkProcessor {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	return &ChunkProcessor{
		chunkSize: chunkSize,
		completed: false,
	}
}

func (cp *ChunkProcessor) ProcessFile(file *multipart.FileHeader, processor func(chunk ChunkData) ChunkResult) ([]ChunkResult, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cp.totalSize = file.Size
	var results []ChunkResult
	chunkIndex := 0
	buffer := make([]byte, cp.chunkSize)

	for {
		n, err := f.Read(buffer)
		if n == 0 {
			break
		}

		isLast := err == io.EOF
		chunkData := ChunkData{
			Index:  chunkIndex,
			Data:   buffer[:n],
			IsLast: isLast,
			Offset: cp.processedSize,
			Size:   n,
		}

		result := processor(chunkData)
		results = append(results, result)

		cp.processedSize += int64(n)
		chunkIndex++

		if err == io.EOF {
			break
		}
		if err != nil {
			return results, err
		}
	}

	cp.completed = true
	return results, nil
}

func (cp *ChunkProcessor) ProcessCSVFile(file *multipart.FileHeader, lineProcessor func(line string, lineNumber int) error) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		if err := lineProcessor(line, lineNumber); err != nil {
			return err
		}
		lineNumber++
	}

	return scanner.Err()
}

func (cp *ChunkProcessor) ProcessCSVInChunks(file *multipart.FileHeader, chunkProcessor func(lines []string, startLine int) error) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	lineNumber := 0
	startLine := 0

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		// Process chunk when it reaches the desired size
		if len(lines) >= cp.getChunkLines() {
			if err := chunkProcessor(lines, startLine); err != nil {
				return err
			}
			startLine = lineNumber + 1
			lines = nil // Reset slice
		}
		lineNumber++
	}

	// Process remaining lines
	if len(lines) > 0 {
		if err := chunkProcessor(lines, startLine); err != nil {
			return err
		}
	}

	cp.completed = true
	return scanner.Err()
}

func (cp *ChunkProcessor) SplitByLines(data []byte) []string {
	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func (cp *ChunkProcessor) GetProgress() float64 {
	if cp.totalSize == 0 {
		return 0
	}
	return float64(cp.processedSize) / float64(cp.totalSize) * 100
}

func (cp *ChunkProcessor) GetChunkSize() int {
	return cp.chunkSize
}

func (cp *ChunkProcessor) IsCompleted() bool {
	return cp.completed
}

func (cp *ChunkProcessor) Reset() {
	cp.processedSize = 0
	cp.totalSize = 0
	cp.completed = false
}

// Calculate approximate number of lines per chunk based on average line length
func (cp *ChunkProcessor) getChunkLines() int {
	// Assuming average line length of 100 characters
	avgLineLength := 100
	return cp.chunkSize / avgLineLength
}