package parsers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// FileFormat representa os formatos de arquivo suportados
type FileFormat int

const (
	FormatUnknown FileFormat = iota
	FormatCSV
	FormatXLSX
	FormatXLS
	FormatXML
	FormatJSON
	FormatFHIRJSON
	FormatFHIRXML
	FormatHL7
	FormatTXT
)

func (f FileFormat) String() string {
	switch f {
	case FormatCSV:
		return "CSV"
	case FormatXLSX:
		return "XLSX"
	case FormatXLS:
		return "XLS"
	case FormatXML:
		return "XML"
	case FormatJSON:
		return "JSON"
	case FormatFHIRJSON:
		return "FHIR_JSON"
	case FormatFHIRXML:
		return "FHIR_XML"
	case FormatHL7:
		return "HL7"
	case FormatTXT:
		return "TXT"
	default:
		return "Unknown"
	}
}

// FormatDetector detecta automaticamente o formato do arquivo
type FormatDetector struct {
	mimeDetector *mimetype.MIME
}

func NewFormatDetector() *FormatDetector {
	return &FormatDetector{}
}

// DetectFormat detecta o formato baseado no conteúdo do arquivo
func (d *FormatDetector) DetectFormat(data []byte, filename string) FileFormat {
	// Usar mimetype para detectar o tipo básico
	mtype := mimetype.Detect(data)
	
	// Primeiro, verificar se é FHIR (priority check)
	if fhirFormat := d.detectFHIRFormat(data); fhirFormat != FormatUnknown {
		return fhirFormat
	}

	// Verificar se é HL7 (priority check)
	if d.isHL7(data) {
		return FormatHL7
	}

	// Mapear MIME types para nossos formatos
	switch mtype.String() {
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return FormatXLSX
	case "application/vnd.ms-excel":
		return FormatXLS
	case "text/csv":
		return FormatCSV
	case "application/xml", "text/xml":
		return FormatXML
	case "application/json":
		return FormatJSON
	case "text/plain":
		// Para text/plain, fazer análise mais detalhada
		return d.detectTextFormat(data, filename)
	default:
		// Fallback para detecção por extensão
		return d.detectByExtension(filename)
	}
}

// detectFHIRFormat detecta especificamente formatos FHIR
func (d *FormatDetector) detectFHIRFormat(data []byte) FileFormat {
	// Tentar JSON FHIR primeiro
	if d.isFHIRJSON(data) {
		return FormatFHIRJSON
	}

	// Tentar XML FHIR
	if d.isFHIRXML(data) {
		return FormatFHIRXML
	}

	return FormatUnknown
}

// isFHIRJSON verifica se é um JSON FHIR válido
func (d *FormatDetector) isFHIRJSON(data []byte) bool {
	// Primeiro, verificar se é JSON válido
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return false
	}

	// Verificar características específicas do FHIR
	resourceType, hasResourceType := jsonData["resourceType"]
	if !hasResourceType {
		return false
	}

	// Lista de resource types FHIR comuns no contexto hospitalar
	fhirResourceTypes := map[string]bool{
		"Patient":       true,
		"Practitioner":  true,
		"Organization":  true,
		"Location":      true,
		"Encounter":     true,
		"Observation":   true,
		"Condition":     true,
		"Procedure":     true,
		"Medication":    true,
		"Bundle":        true,
		"DiagnosticReport": true,
		"AllergyIntolerance": true,
		"Immunization":  true,
		"CarePlan":      true,
		"Device":        true,
		"Appointment":   true,
		"Schedule":      true,
		"Slot":          true,
	}

	resourceTypeStr, ok := resourceType.(string)
	if !ok {
		return false
	}

	// Verificar se é um resource type FHIR válido
	if !fhirResourceTypes[resourceTypeStr] {
		return false
	}

	// Verificações adicionais específicas do FHIR
	// Presença de campos típicos do FHIR
	if resourceTypeStr == "Bundle" {
		// Bundle deve ter 'type' e 'entry'
		_, hasType := jsonData["type"]
		_, hasEntry := jsonData["entry"]
		return hasType || hasEntry
	}

	// Para outros resources, verificar meta ou id
	_, hasMeta := jsonData["meta"]
	_, hasId := jsonData["id"]
	
	return hasMeta || hasId
}

// isFHIRXML verifica se é um XML FHIR válido
func (d *FormatDetector) isFHIRXML(data []byte) bool {
	// Verificar se é XML válido
	decoder := xml.NewDecoder(bytes.NewReader(data))
	
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		if startElement, ok := token.(xml.StartElement); ok {
			// Verificar namespace FHIR
			if startElement.Name.Space == "http://hl7.org/fhir" {
				return true
			}

			// Verificar se o elemento raiz é um resource type FHIR
			fhirElements := []string{
				"Patient", "Practitioner", "Organization", "Location",
				"Encounter", "Observation", "Condition", "Procedure",
				"Medication", "Bundle", "DiagnosticReport", "AllergyIntolerance",
				"Immunization", "CarePlan", "Device", "Appointment", "Schedule", "Slot",
			}

			for _, fhirElement := range fhirElements {
				if startElement.Name.Local == fhirElement {
					return true
				}
			}

			// Se chegou até aqui, não é o primeiro elemento, então para
			break
		}
	}

	return false
}

// isHL7 verifica se é um arquivo HL7 válido
func (d *FormatDetector) isHL7(data []byte) bool {
	content := string(data[:min(2048, len(data))])
	
	// HL7 messages começam com MSH (Message Header)
	if !strings.HasPrefix(content, "MSH") {
		return false
	}
	
	// Verificar estrutura básica de HL7
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return false
	}
	
	// Primeira linha deve ser MSH com separadores válidos
	mshLine := strings.TrimSpace(lines[0])
	if len(mshLine) < 20 {
		return false
	}
	
	// Verificar se tem separadores de campo válidos (|, ^, ~, &)
	validSeparators := []string{"|", "^", "~", "&"}
	separatorCount := 0
	
	for _, sep := range validSeparators {
		if strings.Contains(mshLine, sep) {
			separatorCount++
		}
	}
	
	// Deve ter pelo menos 2 tipos de separadores diferentes
	return separatorCount >= 2
}

// detectTextFormat analisa texto puro para identificar formato
func (d *FormatDetector) detectTextFormat(data []byte, filename string) FileFormat {
	content := string(data[:min(2048, len(data))])
	
	// Verificar se parece CSV
	if d.looksLikeCSV(content) {
		return FormatCSV
	}

	// Verificar por extensão como fallback
	if format := d.detectByExtension(filename); format != FormatUnknown {
		return format
	}

	return FormatTXT
}

// looksLikeCSV verifica se o conteúdo parece CSV
func (d *FormatDetector) looksLikeCSV(content string) bool {
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return false
	}

	// Verificar separadores comuns
	separators := []string{",", ";", "\t", "|"}
	
	for _, sep := range separators {
		if d.hasConsistentSeparator(lines, sep) {
			return true
		}
	}

	return false
}

// hasConsistentSeparator verifica se as linhas têm separador consistente
func (d *FormatDetector) hasConsistentSeparator(lines []string, separator string) bool {
	if len(lines) < 2 {
		return false
	}

	// Contar colunas na primeira linha
	firstLineColumns := strings.Count(lines[0], separator) + 1
	if firstLineColumns < 2 {
		return false
	}

	// Verificar consistência nas próximas linhas
	consistentLines := 0
	for i := 1; i < min(5, len(lines)); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		
		columns := strings.Count(lines[i], separator) + 1
		if abs(columns-firstLineColumns) <= 1 {
			consistentLines++
		}
	}

	return consistentLines >= 2
}

// detectByExtension detecta formato pela extensão do arquivo
func (d *FormatDetector) detectByExtension(filename string) FileFormat {
	if filename == "" {
		return FormatUnknown
	}

	lower := strings.ToLower(filename)
	
	if strings.HasSuffix(lower, ".csv") {
		return FormatCSV
	}
	if strings.HasSuffix(lower, ".xlsx") {
		return FormatXLSX
	}
	if strings.HasSuffix(lower, ".xls") {
		return FormatXLS
	}
	if strings.HasSuffix(lower, ".xml") {
		return FormatXML
	}
	if strings.HasSuffix(lower, ".json") {
		return FormatJSON
	}
	if strings.HasSuffix(lower, ".hl7") {
		return FormatHL7
	}
	if strings.HasSuffix(lower, ".txt") {
		return FormatTXT
	}

	return FormatUnknown
}

// DetectCSVSeparator detecta o separador usado em um CSV
func (d *FormatDetector) DetectCSVSeparator(data []byte) string {
	content := string(data[:min(1024, len(data))])
	lines := strings.Split(content, "\n")
	
	if len(lines) < 2 {
		return ","
	}

	separators := []string{",", ";", "\t", "|"}
	maxCount := 0
	bestSeparator := ","

	for _, sep := range separators {
		if d.hasConsistentSeparator(lines, sep) {
			count := strings.Count(lines[0], sep)
			if count > maxCount {
				maxCount = count
				bestSeparator = sep
			}
		}
	}

	return bestSeparator
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}