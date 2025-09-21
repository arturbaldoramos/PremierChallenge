package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"example.com/m/v2/internal/domain"
	"github.com/google/uuid"
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
	case FormatHL7:
		return p.parseHL7(data)
	case FormatTXT:
		// Tentar como CSV primeiro
		if csvData, _, err := p.parseCSV(data); err == nil {
			return csvData, FormatCSV, nil
		}
		// Tentar como HL7
		if hl7Data, _, err := p.parseHL7(data); err == nil {
			return hl7Data, FormatHL7, nil
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

// parseHL7 processa arquivos HL7
func (p *UnifiedParser) parseHL7(data []byte) (*ParsedData, FileFormat, error) {
	content := string(data)
	lines := strings.Split(content, "\n")
	
	var messages []string
	var currentMessage strings.Builder
	
	// Agrupar linhas em mensagens HL7
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Se começa com MSH, é início de nova mensagem
		if strings.HasPrefix(line, "MSH") {
			if currentMessage.Len() > 0 {
				messages = append(messages, currentMessage.String())
				currentMessage.Reset()
			}
		}
		currentMessage.WriteString(line + "\n")
	}
	
	// Adicionar última mensagem
	if currentMessage.Len() > 0 {
		messages = append(messages, currentMessage.String())
	}
	
	if len(messages) == 0 {
		return nil, FormatHL7, fmt.Errorf("nenhuma mensagem HL7 válida encontrada")
	}
	
	// Headers para HL7
	headers := []string{
		"message_type", "message_control_id", "sending_application", 
		"sending_facility", "receiving_application", "receiving_facility",
		"message_datetime", "security", "message_version", "patient_id",
		"patient_name", "patient_dob", "patient_gender", "patient_address",
		"raw_message",
	}
	
	var rows [][]string
	
	// Processar cada mensagem HL7
	for _, message := range messages {
		row := p.parseHL7Message(message)
		rows = append(rows, row)
	}
	
	metadata := map[string]interface{}{
		"hl7_version": "2.x",
		"total_messages": len(messages),
		"columns": len(headers),
		"original_type": "hl7",
	}
	
	return &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}, FormatHL7, nil
}

// parseHL7Message extrai dados de uma mensagem HL7 individual
func (p *UnifiedParser) parseHL7Message(message string) []string {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return make([]string, 15) // Retornar array vazio com tamanho correto
	}
	
	// Parsear MSH (Message Header)
	mshLine := lines[0]
	mshFields := strings.Split(mshLine, "|")
	
	messageType := ""
	messageControlId := ""
	sendingApp := ""
	sendingFacility := ""
	receivingApp := ""
	receivingFacility := ""
	messageDateTime := ""
	security := ""
	messageVersion := ""
	
	if len(mshFields) >= 12 {
		messageType = mshFields[8]
		messageControlId = mshFields[9]
		sendingApp = mshFields[2]
		sendingFacility = mshFields[3]
		receivingApp = mshFields[4]
		receivingFacility = mshFields[5]
		messageDateTime = mshFields[6]
		security = mshFields[7]
		messageVersion = mshFields[11]
	}
	
	// Parsear PID (Patient Identification) se existir
	patientId := ""
	patientName := ""
	patientDob := ""
	patientGender := ""
	patientAddress := ""
	
	for _, line := range lines {
		if strings.HasPrefix(line, "PID|") {
			pidFields := strings.Split(line, "|")
			if len(pidFields) >= 20 {
				patientId = pidFields[2]
				if len(pidFields) >= 6 {
					patientName = pidFields[5]
				}
				if len(pidFields) >= 8 {
					patientDob = pidFields[7]
				}
				if len(pidFields) >= 9 {
					patientGender = pidFields[8]
				}
				if len(pidFields) >= 12 {
					patientAddress = pidFields[11]
				}
			}
			break
		}
	}
	
	return []string{
		messageType,
		messageControlId,
		sendingApp,
		sendingFacility,
		receivingApp,
		receivingFacility,
		messageDateTime,
		security,
		messageVersion,
		patientId,
		patientName,
		patientDob,
		patientGender,
		patientAddress,
		message, // Mensagem completa
	}
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
	case "estados", "states":
		return p.convertToEstados(data)
	case "cid10", "cid-10":
		return p.convertToCID10(data)
	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}

// convertToHospitals converte para []domain.Hospital
func (p *UnifiedParser) convertToHospitals(data *ParsedData) ([]domain.Hospital, error) {
	var hospitals []domain.Hospital
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range data.Headers {
		headerMap[strings.ToLower(header)] = i
	}
	
	for _, row := range data.Rows {
		if len(row) == 0 {
			continue
		}
		
		hospital := domain.Hospital{}
		
		// UUID
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			if id, err := uuid.Parse(row[idx]); err == nil {
				hospital.UUID = id
			} else {
				hospital.UUID = uuid.New()
			}
		} else {
			hospital.UUID = uuid.New()
		}
		
		// Nome
		if idx, exists := headerMap["nome"]; exists && idx < len(row) {
			hospital.Nome = strings.TrimSpace(row[idx])
		}
		
		// CEP
		if idx, exists := headerMap["cep"]; exists && idx < len(row) {
			hospital.CEP = strings.TrimSpace(row[idx])
		}
		
		// Especialidades
		if idx, exists := headerMap["especialidades"]; exists && idx < len(row) {
			hospital.Especialidades = strings.TrimSpace(row[idx])
		}
		
		// Leitos totais
		if idx, exists := headerMap["leitos_totais"]; exists && idx < len(row) {
			if leitos, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				hospital.LeitosTotais = leitos
			}
		}
		
		// Código do município
		if idx, exists := headerMap["cidade"]; exists && idx < len(row) {
			hospital.CodMunicipio = strings.TrimSpace(row[idx])
		}
		if idx, exists := headerMap["cod_municipio"]; exists && idx < len(row) {
			hospital.CodMunicipio = strings.TrimSpace(row[idx])
		}
		
		// Bairro
		if idx, exists := headerMap["bairro"]; exists && idx < len(row) {
			hospital.Bairro = strings.TrimSpace(row[idx])
		}
		
		// Só adicionar se tiver pelo menos nome
		if hospital.Nome != "" {
			hospitals = append(hospitals, hospital)
		}
	}
	
	return hospitals, nil
}

// convertToPatients converte para []domain.Paciente  
func (p *UnifiedParser) convertToPatients(data *ParsedData) ([]domain.Paciente, error) {
	var patients []domain.Paciente
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range data.Headers {
		headerMap[strings.ToLower(header)] = i
	}
	
	for _, row := range data.Rows {
		if len(row) == 0 {
			continue
		}
		
		paciente := domain.Paciente{}
		
		// ID (UUID do código)
		if idx, exists := headerMap["id"]; exists && idx < len(row) {
			if id, err := uuid.Parse(row[idx]); err == nil {
				paciente.ID = id
			} else {
				paciente.ID = uuid.New()
			}
		} else if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			if id, err := uuid.Parse(row[idx]); err == nil {
				paciente.ID = id
			} else {
				paciente.ID = uuid.New()
			}
		} else {
			paciente.ID = uuid.New()
		}

		// CPF
		if idx, exists := headerMap["cpf"]; exists && idx < len(row) {
			paciente.CPF = strings.TrimSpace(row[idx])
		}

		// Nome
		if idx, exists := headerMap["nome"]; exists && idx < len(row) {
			paciente.Nome = strings.TrimSpace(row[idx])
		}
		if idx, exists := headerMap["nome_completo"]; exists && idx < len(row) {
			paciente.Nome = strings.TrimSpace(row[idx])
		}

		// Gênero
		if idx, exists := headerMap["genero"]; exists && idx < len(row) {
			paciente.Genero = strings.TrimSpace(row[idx])
		}

		// Código do município
		if idx, exists := headerMap["cod_municipio"]; exists && idx < len(row) {
			paciente.CodMunicipio = strings.TrimSpace(row[idx])
		}

		// Bairro
		if idx, exists := headerMap["bairro"]; exists && idx < len(row) {
			paciente.Bairro = strings.TrimSpace(row[idx])
		}

		// Convênio
		if idx, exists := headerMap["convenio"]; exists && idx < len(row) {
			paciente.Convenio = strings.TrimSpace(row[idx])
		}

		// CID10
		if idx, exists := headerMap["cid10"]; exists && idx < len(row) {
			paciente.CID10 = strings.TrimSpace(row[idx])
		}
		if idx, exists := headerMap["cid-10"]; exists && idx < len(row) {
			paciente.CID10 = strings.TrimSpace(row[idx])
		}

		// Só adicionar se tiver pelo menos CPF ou nome
		if paciente.CPF != "" || paciente.Nome != "" {
			patients = append(patients, paciente)
		}
	}
	
	return patients, nil
}

// convertToMedicos converte para []domain.Medico
func (p *UnifiedParser) convertToMedicos(data *ParsedData) ([]domain.Medico, error) {
	var medicos []domain.Medico
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range data.Headers {
		headerMap[strings.ToLower(header)] = i
	}
	
	for _, row := range data.Rows {
		if len(row) == 0 {
			continue
		}
		
		medico := domain.Medico{}
		
		// UUID
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			if id, err := uuid.Parse(row[idx]); err == nil {
				medico.UUID = id
			} else {
				medico.UUID = uuid.New()
			}
		} else {
			medico.UUID = uuid.New()
		}
		
		// Nome
		if idx, exists := headerMap["nome"]; exists && idx < len(row) {
			medico.Nome = strings.TrimSpace(row[idx])
		}
		if idx, exists := headerMap["nome_completo"]; exists && idx < len(row) {
			medico.Nome = strings.TrimSpace(row[idx])
		}
		
		// Especialidade
		if idx, exists := headerMap["especialidade"]; exists && idx < len(row) {
			medico.Especialidade = strings.TrimSpace(row[idx])
		}
		
		// Código do município
		if idx, exists := headerMap["cod_municipio"]; exists && idx < len(row) {
			medico.CodMunicipio = strings.TrimSpace(row[idx])
		}
		if idx, exists := headerMap["cidade"]; exists && idx < len(row) {
			medico.CodMunicipio = strings.TrimSpace(row[idx])
		}
		
		// Só adicionar se tiver pelo menos nome
		if medico.Nome != "" {
			medicos = append(medicos, medico)
		}
	}
	
	return medicos, nil
}

// convertToMunicipios converte para []domain.Municipio
func (p *UnifiedParser) convertToMunicipios(data *ParsedData) ([]domain.Municipio, error) {
	var municipios []domain.Municipio
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range data.Headers {
		headerMap[strings.ToLower(header)] = i
	}
	
	for _, row := range data.Rows {
		if len(row) == 0 {
			continue
		}
		
		municipio := domain.Municipio{}
		
		// Código IBGE
		if idx, exists := headerMap["codigo_ibge"]; exists && idx < len(row) {
			municipio.Codigo = strings.TrimSpace(row[idx])
		}
		
		// Nome
		if idx, exists := headerMap["nome"]; exists && idx < len(row) {
			municipio.Nome = strings.TrimSpace(row[idx])
		}
		
		// Latitude
		if idx, exists := headerMap["latitude"]; exists && idx < len(row) {
			municipio.Latitude = strings.TrimSpace(row[idx])
		}
		
		// Longitude
		if idx, exists := headerMap["longitude"]; exists && idx < len(row) {
			municipio.Longitude = strings.TrimSpace(row[idx])
		}
		
		// Capital
		if idx, exists := headerMap["capital"]; exists && idx < len(row) {
			municipio.Capital = strings.TrimSpace(row[idx])
		}
		
		// Código UF
		if idx, exists := headerMap["codigo_uf"]; exists && idx < len(row) {
			municipio.CodigoUF = strings.TrimSpace(row[idx])
		}
		
		// SIAFI ID
		if idx, exists := headerMap["siafi_id"]; exists && idx < len(row) {
			municipio.SiafiId = strings.TrimSpace(row[idx])
		}
		
		// DDD
		if idx, exists := headerMap["ddd"]; exists && idx < len(row) {
			municipio.DDD = strings.TrimSpace(row[idx])
		}
		
		// Fuso horário
		if idx, exists := headerMap["fuso_horario"]; exists && idx < len(row) {
			municipio.FusoHora = strings.TrimSpace(row[idx])
		}
		
		// População
		if idx, exists := headerMap["populacao"]; exists && idx < len(row) {
			if pop, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
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

// convertToEstados converte para []domain.Estado
func (p *UnifiedParser) convertToEstados(data *ParsedData) ([]domain.Estado, error) {
	var estados []domain.Estado
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range data.Headers {
		headerMap[strings.ToLower(header)] = i
	}
	
	for _, row := range data.Rows {
		if len(row) == 0 {
			continue
		}
		
		estado := domain.Estado{}
		
		// Código UF
		if idx, exists := headerMap["codigo_uf"]; exists && idx < len(row) {
			estado.Codigo = strings.TrimSpace(row[idx])
		}
		
		// UF
		if idx, exists := headerMap["uf"]; exists && idx < len(row) {
			estado.UnidadeFederativa = strings.TrimSpace(row[idx])
		}
		
		// Nome
		if idx, exists := headerMap["nome"]; exists && idx < len(row) {
			estado.Nome = strings.TrimSpace(row[idx])
		}
		
		// Região
		if idx, exists := headerMap["regiao"]; exists && idx < len(row) {
			estado.Regiao = strings.TrimSpace(row[idx])
		}
		
		// Latitude
		if idx, exists := headerMap["latitude"]; exists && idx < len(row) {
			estado.Latitude = strings.TrimSpace(row[idx])
		}
		
		// Longitude
		if idx, exists := headerMap["longitude"]; exists && idx < len(row) {
			estado.Longitude = strings.TrimSpace(row[idx])
		}
		
		// Só adicionar se tiver pelo menos código e nome
		if estado.Codigo != "" && estado.Nome != "" {
			estados = append(estados, estado)
		}
	}
	
	return estados, nil
}

// convertToCID10 converte para []domain.Cid10
func (p *UnifiedParser) convertToCID10(data *ParsedData) ([]domain.Cid10, error) {
	var cid10s []domain.Cid10
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range data.Headers {
		headerMap[strings.ToLower(header)] = i
	}
	
	for _, row := range data.Rows {
		if len(row) == 0 {
			continue
		}
		
		cid10 := domain.Cid10{}
		
		// ID
		if idx, exists := headerMap["id"]; exists && idx < len(row) {
			if id, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				cid10.ID = id
			}
		}
		
		// Código
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
			cid10.Codigo = strings.TrimSpace(row[idx])
		}
		
		// Descrição
		if idx, exists := headerMap["descricao"]; exists && idx < len(row) {
			cid10.Descricao = strings.TrimSpace(row[idx])
		}
		
		// Categoria
		if idx, exists := headerMap["categoria"]; exists && idx < len(row) {
			cid10.Categoria = strings.TrimSpace(row[idx])
		}
		
		// Grupo
		if idx, exists := headerMap["grupo"]; exists && idx < len(row) {
			cid10.Grupo = strings.TrimSpace(row[idx])
		}
		
		// Só adicionar se tiver pelo menos código e descrição
		if cid10.Codigo != "" && cid10.Descricao != "" {
			cid10s = append(cid10s, cid10)
		}
	}
	
	return cid10s, nil
}