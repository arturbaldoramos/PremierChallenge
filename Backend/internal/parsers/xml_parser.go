package parsers

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"example.com/m/v2/internal/domain"
	"github.com/google/uuid"
)

// XMLPacientesRoot representa a estrutura raiz do XML de pacientes
type XMLPacientesRoot struct {
	XMLName   xml.Name     `xml:"Pacientes"`
	Pacientes []XMLPaciente `xml:"Paciente"`
}

// XMLPaciente representa um paciente no XML
type XMLPaciente struct {
	Codigo       string `xml:"Codigo"`
	CPF          string `xml:"CPF"`
	NomeCompleto string `xml:"Nome_Completo"`
	Genero       string `xml:"Genero"`
	CodMunicipio string `xml:"Cod_municipio"`
	Bairro       string `xml:"Bairro"`
	Convenio     string `xml:"Convenio"`
	CID10        string `xml:"CID-10"`
}

// XMLHospitaisRoot representa a estrutura raiz do XML de hospitais
type XMLHospitaisRoot struct {
	XMLName   xml.Name      `xml:"Hospitais"`
	Hospitais []XMLHospital `xml:"Hospital"`
}

// XMLHospital representa um hospital no XML
type XMLHospital struct {
	Codigo         string `xml:"Codigo"`
	Nome           string `xml:"Nome"`
	CEP            string `xml:"CEP"`
	Especialidades string `xml:"Especialidades"`
	LeitosTotais   string `xml:"Leitos_Totais"`
	CodMunicipio   string `xml:"Cod_municipio"`
	Bairro         string `xml:"Bairro"`
}

// XMLMedicosRoot representa a estrutura raiz do XML de médicos
type XMLMedicosRoot struct {
	XMLName xml.Name    `xml:"Medicos"`
	Medicos []XMLMedico `xml:"Medico"`
}

// XMLMedico representa um médico no XML
type XMLMedico struct {
	Codigo        string `xml:"Codigo"`
	NomeCompleto  string `xml:"Nome_Completo"`
	Especialidade string `xml:"Especialidade"`
	CodMunicipio  string `xml:"Cod_municipio"`
}

// XMLParser parser especializado para XML estruturado
type XMLParser struct{}

func NewXMLParser() *XMLParser {
	return &XMLParser{}
}

// ParseXMLData detecta e parseia diferentes tipos de XML estruturado
func (p *XMLParser) ParseXMLData(data []byte) (*ParsedData, string, error) {
	// Tentar detectar o tipo de XML pela estrutura
	content := string(data[:min(1000, len(data))]) // Verificar apenas primeiros 1000 chars para performance

	if strings.Contains(content, "<Pacientes>") || strings.Contains(content, "<Paciente>") {
		fmt.Printf("🔍 XML: Detectado como XML de pacientes\n")
		return p.parsePacientesXML(data)
	}

	if strings.Contains(content, "<Hospitais>") || strings.Contains(content, "<Hospital>") {
		fmt.Printf("🔍 XML: Detectado como XML de hospitais\n")
		return p.parseHospitaisXML(data)
	}

	if strings.Contains(content, "<Medicos>") || strings.Contains(content, "<Medico>") {
		fmt.Printf("🔍 XML: Detectado como XML de médicos\n")
		return p.parseMedicosXML(data)
	}

	fmt.Printf("🔍 XML: Não detectou tipo específico, usando parser genérico\n")
	// XML genérico - tentar estrutura básica
	return p.parseGenericXML(data)
}

// parsePacientesXML parseia XML de pacientes
func (p *XMLParser) parsePacientesXML(data []byte) (*ParsedData, string, error) {
	fmt.Printf("🔄 XML: Iniciando parse de pacientes (tamanho: %d bytes)\n", len(data))

	var root XMLPacientesRoot

	if err := xml.Unmarshal(data, &root); err != nil {
		fmt.Printf("❌ XML: Erro no Unmarshal: %v\n", err)
		return nil, "", fmt.Errorf("erro ao fazer parse do XML de pacientes: %v", err)
	}

	fmt.Printf("✅ XML: Parse bem-sucedido, %d pacientes encontrados\n", len(root.Pacientes))
	
	// Headers para pacientes
	headers := []string{
		"codigo", "cpf", "nome_completo", "genero",
		"cod_municipio", "bairro", "convenio", "cid10",
	}
	
	var rows [][]string
	
	for _, paciente := range root.Pacientes {
		row := []string{
			paciente.Codigo,
			paciente.CPF,
			paciente.NomeCompleto,
			paciente.Genero,
			paciente.CodMunicipio,
			paciente.Bairro,
			paciente.Convenio,
			paciente.CID10,
		}
		rows = append(rows, row)
	}
	
	metadata := map[string]interface{}{
		"xml_type": "pacientes",
		"total_records": len(root.Pacientes),
		"columns": len(headers),
	}
	
	parsedData := &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}
	
	return parsedData, "pacientes", nil
}

// parseHospitaisXML parseia XML de hospitais
func (p *XMLParser) parseHospitaisXML(data []byte) (*ParsedData, string, error) {
	var root XMLHospitaisRoot
	
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, "", fmt.Errorf("erro ao fazer parse do XML de hospitais: %v", err)
	}
	
	headers := []string{
		"codigo", "nome", "cep", "especialidades", 
		"leitos_totais", "cod_municipio", "bairro",
	}
	
	var rows [][]string
	
	for _, hospital := range root.Hospitais {
		row := []string{
			hospital.Codigo,
			hospital.Nome,
			hospital.CEP,
			hospital.Especialidades,
			hospital.LeitosTotais,
			hospital.CodMunicipio,
			hospital.Bairro,
		}
		rows = append(rows, row)
	}
	
	metadata := map[string]interface{}{
		"xml_type": "hospitais",
		"total_records": len(root.Hospitais),
		"columns": len(headers),
	}
	
	parsedData := &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}
	
	return parsedData, "hospitais", nil
}

// parseMedicosXML parseia XML de médicos
func (p *XMLParser) parseMedicosXML(data []byte) (*ParsedData, string, error) {
	var root XMLMedicosRoot
	
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, "", fmt.Errorf("erro ao fazer parse do XML de médicos: %v", err)
	}
	
	headers := []string{
		"codigo", "nome_completo", "especialidade", "cod_municipio",
	}
	
	var rows [][]string
	
	for _, medico := range root.Medicos {
		row := []string{
			medico.Codigo,
			medico.NomeCompleto,
			medico.Especialidade,
			medico.CodMunicipio,
		}
		rows = append(rows, row)
	}
	
	metadata := map[string]interface{}{
		"xml_type": "medicos",
		"total_records": len(root.Medicos),
		"columns": len(headers),
	}
	
	parsedData := &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}
	
	return parsedData, "medicos", nil
}

// parseGenericXML parseia XML genérico
func (p *XMLParser) parseGenericXML(data []byte) (*ParsedData, string, error) {
	// Parser XML genérico básico
	headers := []string{"element", "value"}
	rows := [][]string{}
	
	metadata := map[string]interface{}{
		"xml_type": "generic",
		"note": "Generic XML parsing - limited structure detection",
	}
	
	parsedData := &ParsedData{
		Headers:  headers,
		Rows:     rows,
		Metadata: metadata,
	}
	
	return parsedData, "unknown", nil
}

// ConvertXMLToDomain converte dados XML parseados para domain objects
func (p *XMLParser) ConvertXMLToDomain(parsedData *ParsedData, dataType string) (interface{}, error) {
	switch strings.ToLower(dataType) {
	case "pacientes":
		return p.convertToPacientes(parsedData)
	case "hospitais":
		return p.convertToHospitais(parsedData)
	case "medicos":
		return p.convertToMedicos(parsedData)
	default:
		return nil, fmt.Errorf("tipo de dados não suportado: %s", dataType)
	}
}

// convertToPacientes converte dados XML para []domain.Paciente
func (p *XMLParser) convertToPacientes(parsedData *ParsedData) ([]domain.Paciente, error) {
	var pacientes []domain.Paciente
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		headerMap[header] = i
	}
	
	for _, row := range parsedData.Rows {
		if len(row) == 0 {
			continue
		}
		
		paciente := domain.Paciente{}
		
		// ID (UUID do código)
		if idx, exists := headerMap["codigo"]; exists && idx < len(row) {
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

		// Só adicionar se tiver pelo menos CPF
		if paciente.CPF != "" {
			pacientes = append(pacientes, paciente)
		}
	}
	
	return pacientes, nil
}

// convertToHospitais converte dados XML para []domain.Hospital
func (p *XMLParser) convertToHospitais(parsedData *ParsedData) ([]domain.Hospital, error) {
	var hospitais []domain.Hospital
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		headerMap[header] = i
	}
	
	for _, row := range parsedData.Rows {
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
		if idx, exists := headerMap["cod_municipio"]; exists && idx < len(row) {
			hospital.CodMunicipio = strings.TrimSpace(row[idx])
		}
		
		// Bairro
		if idx, exists := headerMap["bairro"]; exists && idx < len(row) {
			hospital.Bairro = strings.TrimSpace(row[idx])
		}
		
		// Só adicionar se tiver pelo menos nome
		if hospital.Nome != "" {
			hospitais = append(hospitais, hospital)
		}
	}
	
	return hospitais, nil
}

// convertToMedicos converte dados XML para []domain.Medico
func (p *XMLParser) convertToMedicos(parsedData *ParsedData) ([]domain.Medico, error) {
	var medicos []domain.Medico
	
	// Mapear headers para índices
	headerMap := make(map[string]int)
	for i, header := range parsedData.Headers {
		headerMap[header] = i
	}
	
	for _, row := range parsedData.Rows {
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
		}
		
		// Nome
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
		
		// Só adicionar se tiver pelo menos nome
		if medico.Nome != "" {
			medicos = append(medicos, medico)
		}
	}
	
	return medicos, nil
}