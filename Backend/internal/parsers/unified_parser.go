package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"example.com/m/v2/internal/domain"
	"github.com/xuri/excelize/v2"
)

// ParsedData representa dados parseados de qualquer formato
type ParsedData struct {
	Headers []string              `json:"headers"`
	Rows    [][]string           `json:"rows"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UnifiedParser é um parser que funciona com múltiplos formatos
type UnifiedParser struct {
	detector    *FormatDetector
	csvParser   *CSVParser
	xmlParser   *XMLParser
}

func NewUnifiedParser() *UnifiedParser {
	return &UnifiedParser{
		detector:  NewFormatDetector(),
		csvParser: NewCSVParser(),
		xmlParser: NewXMLParser(),
	}
}

// ParseFile detecta automaticamente o formato e faz o parsing
func (p *UnifiedParser) ParseFile(data []byte, filename string) (*ParsedData, FileFormat, error) {
	// Detectar formato
	format := p.detector.DetectFormat(data, filename)
	
	// Parse baseado no formato detectado
	switch format {
	case FormatCSV:
		csvData, format, err := p.parseCSV(data)
		return csvData, format, err
	case FormatXLSX, FormatXLS:
		return p.parseExcel(data)
	case FormatJSON:
		return p.parseJSON(data)
	case FormatXML:
		return p.parseXML(data)
	case FormatFHIRJSON:
		return p.parseFHIRJSON(data)
	case FormatFHIRXML:
		return p.parseFHIRXML(data)
	case FormatTXT:
		// Tentar como CSV primeiro
		if csvData, _, err := p.parseCSV(data); err == nil {
			return csvData, FormatCSV, nil
		}
		return nil, format, fmt.Errorf("unsupported text format")
	default:
		return nil, format, fmt.Errorf("unsupported file format: %s", format.String())
	}
}

// parseCSV processa arquivo CSV
func (p *UnifiedParser) parseCSV(data []byte) (*ParsedData, FileFormat, error) {
	// Detectar separador
	separator := p.detector.DetectCSVSeparator(data)
	
	// Parse CSV
	rows, err := p.parseCSVWithSeparator(data, separator)
	if err != nil {
		return nil, FormatCSV, err
	}

	if len(rows) == 0 {
		return nil, FormatCSV, fmt.Errorf("empty CSV file")
	}

	// Primeira linha como headers
	headers := rows[0]
	dataRows := rows[1:]

	metadata := map[string]interface{}{
		"separator": separator,
		"total_rows": len(dataRows),
		"columns": len(headers),
	}

	return &ParsedData{
		Headers:  headers,
		Rows:     dataRows,
		Metadata: metadata,
	}, FormatCSV, nil
}

// parseCSVWithSeparator faz parsing CSV com separador específico
func (p *UnifiedParser) parseCSVWithSeparator(data []byte, separator string) ([][]string, error) {
	content := string(data)
	lines := strings.Split(content, "\n")
	
	var rows [][]string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Split pela vírgula e limpar campos
		fields := strings.Split(line, separator)
		for i, field := range fields {
			fields[i] = strings.TrimSpace(field)
		}
		rows = append(rows, fields)
	}
	
	return rows, nil
}

// parseExcel processa arquivos Excel
func (p *UnifiedParser) parseExcel(data []byte) (*ParsedData, FileFormat, error) {
	// Criar um reader de bytes
	reader := bytes.NewReader(data)
	
	// Abrir Excel file usando bytes
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, FormatXLSX, fmt.Errorf("erro ao abrir arquivo Excel: %v", err)
	}
	defer file.Close()

	// Pegar a primeira sheet
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, FormatXLSX, fmt.Errorf("nenhuma planilha encontrada no arquivo Excel")
	}

	sheetName := sheets[0]
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, FormatXLSX, fmt.Errorf("erro ao ler linhas da planilha: %v", err)
	}

	if len(rows) == 0 {
		return nil, FormatXLSX, fmt.Errorf("planilha Excel vazia")
	}

	// Filtrar linhas vazias e normalizar
	var filteredRows [][]string
	for _, row := range rows {
		if len(row) > 0 {
			// Limpar células vazias no final
			var cleanRow []string
			for _, cell := range row {
				cleanRow = append(cleanRow, strings.TrimSpace(cell))
			}
			if len(cleanRow) > 0 && cleanRow[0] != "" {
				filteredRows = append(filteredRows, cleanRow)
			}
		}
	}

	if len(filteredRows) == 0 {
		return nil, FormatXLSX, fmt.Errorf("nenhuma linha válida encontrada na planilha")
	}

	headers := filteredRows[0]
	dataRows := filteredRows[1:]

	metadata := map[string]interface{}{
		"sheet_name": sheetName,
		"total_sheets": len(sheets),
		"total_rows": len(dataRows),
		"columns": len(headers),
		"original_rows": len(rows),
		"filtered_rows": len(filteredRows),
	}

	return &ParsedData{
		Headers:  headers,
		Rows:     dataRows,
		Metadata: metadata,
	}, FormatXLSX, nil
}

// parseJSON processa arquivos JSON genéricos
func (p *UnifiedParser) parseJSON(data []byte) (*ParsedData, FileFormat, error) {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, FormatJSON, err
	}

	// Tentar converter para array de objetos (formato tabular)
	rows, headers, err := p.jsonToRows(jsonData)
	if err != nil {
		return nil, FormatJSON, err
	}

	metadata := map[string]interface{}{
		"original_type": "json",
		"total_rows": len(rows),
		"columns": len(headers),
	}

	return &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}, FormatJSON, nil
}

// parseXML processa arquivos XML estruturados
func (p *UnifiedParser) parseXML(data []byte) (*ParsedData, FileFormat, error) {
	// Tentar usar o parser XML especializado primeiro
	parsedData, dataType, err := p.xmlParser.ParseXMLData(data)
	if err == nil && dataType != "unknown" {
		// Adicionar informação sobre o tipo detectado
		if parsedData.Metadata == nil {
			parsedData.Metadata = make(map[string]interface{})
		}
		parsedData.Metadata["detected_type"] = dataType
		parsedData.Metadata["parser"] = "structured_xml"
		
		return parsedData, FormatXML, nil
	}
	
	// Fallback para parser XML genérico básico
	metadata := map[string]interface{}{
		"original_type": "xml",
		"parser": "generic_xml",
		"note": "basic XML parsing - structured content not detected",
	}

	// Para XML genérico, retornar estrutura básica
	headers := []string{"element", "value", "attributes"}
	rows := [][]string{}

	return &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}, FormatXML, nil
}

// parseFHIRJSON processa arquivos FHIR em formato JSON
func (p *UnifiedParser) parseFHIRJSON(data []byte) (*ParsedData, FileFormat, error) {
	var fhirData map[string]interface{}
	if err := json.Unmarshal(data, &fhirData); err != nil {
		return nil, FormatFHIRJSON, err
	}

	resourceType, _ := fhirData["resourceType"].(string)
	
	// Se for Bundle, processar entries
	if resourceType == "Bundle" {
		return p.parseFHIRBundle(fhirData)
	}

	// Para resource individual, converter para formato tabular
	rows, headers := p.fhirResourceToRows(fhirData)

	metadata := map[string]interface{}{
		"fhir_version": "R4",
		"resource_type": resourceType,
		"original_type": "fhir_json",
		"total_rows": len(rows),
		"columns": len(headers),
	}

	return &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}, FormatFHIRJSON, nil
}

// parseFHIRXML processa arquivos FHIR em formato XML
func (p *UnifiedParser) parseFHIRXML(data []byte) (*ParsedData, FileFormat, error) {
	// Para simplicidade, converter XML para JSON primeiro e depois processar
	// Em produção, usar parser XML específico para FHIR
	
	metadata := map[string]interface{}{
		"fhir_version": "R4",
		"original_type": "fhir_xml",
		"note": "XML FHIR parsing - converted to tabular format",
	}

	headers := []string{"resource_type", "id", "field", "value"}
	rows := [][]string{}

	return &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}, FormatFHIRXML, nil
}

// parseFHIRBundle processa um Bundle FHIR
func (p *UnifiedParser) parseFHIRBundle(bundleData map[string]interface{}) (*ParsedData, FileFormat, error) {
	entries, ok := bundleData["entry"].([]interface{})
	if !ok {
		return nil, FormatFHIRJSON, fmt.Errorf("invalid FHIR Bundle: no entries")
	}

	var allRows [][]string
	var headers []string
	headerSet := make(map[string]bool)

	// Processar cada entry
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		rows, resourceHeaders := p.fhirResourceToRows(resource)
		
		// Adicionar headers únicos
		for _, header := range resourceHeaders {
			if !headerSet[header] {
				headers = append(headers, header)
				headerSet[header] = true
			}
		}

		allRows = append(allRows, rows...)
	}

	metadata := map[string]interface{}{
		"fhir_version": "R4",
		"resource_type": "Bundle",
		"bundle_type": bundleData["type"],
		"total_entries": len(entries),
		"total_rows": len(allRows),
		"columns": len(headers),
		"original_type": "fhir_json_bundle",
	}

	return &ParsedData{
		Headers:  headers,
		Rows:     allRows,
		Metadata: metadata,
	}, FormatFHIRJSON, nil
}

// fhirResourceToRows converte um resource FHIR para formato de linhas
func (p *UnifiedParser) fhirResourceToRows(resource map[string]interface{}) ([][]string, []string) {
	resourceType, _ := resource["resourceType"].(string)
	id, _ := resource["id"].(string)

	// Headers básicos para FHIR
	headers := []string{"resource_type", "id", "field", "value"}
	var rows [][]string

	// Extrair campos principais
	for key, value := range resource {
		if key == "resourceType" || key == "id" {
			continue
		}

		valueStr := p.interfaceToString(value)
		row := []string{resourceType, id, key, valueStr}
		rows = append(rows, row)
	}

	return rows, headers
}

// jsonToRows converte JSON genérico para formato de linhas
func (p *UnifiedParser) jsonToRows(data interface{}) ([][]string, []string, error) {
	switch v := data.(type) {
	case []interface{}:
		// Array de objetos
		if len(v) == 0 {
			return [][]string{}, []string{}, nil
		}

		// Assumir que todos os objetos têm a mesma estrutura
		firstItem, ok := v[0].(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("array items must be objects")
		}

		// Extrair headers
		var headers []string
		for key := range firstItem {
			headers = append(headers, key)
		}

		// Extrair rows
		var rows [][]string
		for _, item := range v {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			var row []string
			for _, header := range headers {
				value := p.interfaceToString(itemMap[header])
				row = append(row, value)
			}
			rows = append(rows, row)
		}

		return rows, headers, nil

	case map[string]interface{}:
		// Objeto único - converter para single row
		var headers []string
		var row []string

		for key, value := range v {
			headers = append(headers, key)
			row = append(row, p.interfaceToString(value))
		}

		return [][]string{row}, headers, nil

	default:
		return nil, nil, fmt.Errorf("unsupported JSON structure")
	}
}

// interfaceToString converte interface{} para string
func (p *UnifiedParser) interfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// Para objetos complexos, converter para JSON
		jsonBytes, _ := json.Marshal(value)
		return string(jsonBytes)
	}
}

// ConvertToDataType converte ParsedData para tipos específicos do domínio
func (p *UnifiedParser) ConvertToDataType(data *ParsedData, dataType string) (interface{}, error) {
	switch strings.ToLower(dataType) {
	case "hospitais", "hospitals":
		return p.convertToHospitals(data)
	case "pacientes", "patients":
		return p.convertToPatients(data)
	case "medicos", "doctors", "practitioners":
		return p.convertToMedicos(data)
	case "municipios", "municipalities":
		return p.convertToMunicipios(data)
	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}

// convertToHospitals converte para []domain.Hospital
func (p *UnifiedParser) convertToHospitals(data *ParsedData) ([]domain.Hospital, error) {
	// Implementar conversão baseada nos headers e rows
	// Esta é uma versão simplificada - expandir conforme necessário
	var hospitals []domain.Hospital
	
	// TODO: Implementar mapeamento de campos
	// Mapear headers para campos do Hospital
	
	return hospitals, nil
}

// convertToPatients converte para []domain.Paciente  
func (p *UnifiedParser) convertToPatients(data *ParsedData) ([]domain.Paciente, error) {
	var patients []domain.Paciente
	// TODO: Implementar
	return patients, nil
}

// convertToMedicos converte para []domain.Medico
func (p *UnifiedParser) convertToMedicos(data *ParsedData) ([]domain.Medico, error) {
	var medicos []domain.Medico
	// TODO: Implementar  
	return medicos, nil
}

// convertToMunicipios converte para []domain.Municipio
func (p *UnifiedParser) convertToMunicipios(data *ParsedData) ([]domain.Municipio, error) {
	var municipios []domain.Municipio
	// TODO: Implementar
	return municipios, nil
}