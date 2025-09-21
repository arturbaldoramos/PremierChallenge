package parsers

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	RG           string `xml:"RG"`
	DataNascimento string `xml:"Data_Nascimento"`
	Genero       string `xml:"Genero"`
	TipoSanguineo string `xml:"Tipo_Sanguineo"`
	Endereco     string `xml:"Endereco"`
	CodMunicipio string `xml:"Cod_municipio"`
	CEP          string `xml:"CEP"`
	Telefone     string `xml:"Telefone"`
	Email        string `xml:"Email"`
	ContatoEmergencia string `xml:"Contato_Emergencia"`
	Convenio     string `xml:"Convenio"`
	NumeroCarteira string `xml:"Numero_Carteira"`
	Status       string `xml:"Status"`
	Bairro       string `xml:"Bairro"`
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
	content := string(data)
	
	if strings.Contains(content, "<Pacientes>") || strings.Contains(content, "<Paciente>") {
		return p.parsePacientesXML(data)
	}
	
	if strings.Contains(content, "<Hospitais>") || strings.Contains(content, "<Hospital>") {
		return p.parseHospitaisXML(data)
	}
	
	if strings.Contains(content, "<Medicos>") || strings.Contains(content, "<Medico>") {
		return p.parseMedicosXML(data)
	}
	
	// XML genérico - tentar estrutura básica
	return p.parseGenericXML(data)
}

// parsePacientesXML parseia XML de pacientes
func (p *XMLParser) parsePacientesXML(data []byte) (*ParsedData, string, error) {
	var root XMLPacientesRoot
	
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, "", fmt.Errorf("erro ao fazer parse do XML de pacientes: %v", err)
	}
	
	// Headers para pacientes
	headers := []string{
		"codigo", "cpf", "nome_completo", "rg", "data_nascimento", 
		"genero", "tipo_sanguineo", "endereco", "cod_municipio", 
		"cep", "telefone", "email", "contato_emergencia", 
		"convenio", "numero_carteira", "status", "bairro", "cid10",
	}
	
	var rows [][]string
	
	for _, paciente := range root.Pacientes {
		row := []string{
			paciente.Codigo,
			paciente.CPF,
			paciente.NomeCompleto,
			paciente.RG,
			paciente.DataNascimento,
			paciente.Genero,
			paciente.TipoSanguineo,
			paciente.Endereco,
			paciente.CodMunicipio,
			paciente.CEP,
			paciente.Telefone,
			paciente.Email,
			paciente.ContatoEmergencia,
			paciente.Convenio,
			paciente.NumeroCarteira,
			paciente.Status,
			paciente.Bairro,
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
		}
		
		// CPF
		if idx, exists := headerMap["cpf"]; exists && idx < len(row) {
			paciente.CPF = strings.TrimSpace(row[idx])
		}
		
		// Nome
		if idx, exists := headerMap["nome_completo"]; exists && idx < len(row) {
			paciente.Nome = strings.TrimSpace(row[idx])
		}
		
		// RG
		if idx, exists := headerMap["rg"]; exists && idx < len(row) {
			paciente.RG = strings.TrimSpace(row[idx])
		}
		
		// Data de nascimento
		if idx, exists := headerMap["data_nascimento"]; exists && idx < len(row) {
			if dateStr := strings.TrimSpace(row[idx]); dateStr != "" {
				// Tentar diferentes formatos de data
				formats := []string{"2006-01-02", "02/01/2006", "02-01-2006"}
				for _, format := range formats {
					if date, err := time.Parse(format, dateStr); err == nil {
						paciente.DataNascimento = date
						break
					}
				}
			}
		}
		
		// Gênero
		if idx, exists := headerMap["genero"]; exists && idx < len(row) {
			paciente.Genero = strings.TrimSpace(row[idx])
		}
		
		// Tipo sanguíneo
		if idx, exists := headerMap["tipo_sanguineo"]; exists && idx < len(row) {
			paciente.TipoSanguineo = strings.TrimSpace(row[idx])
		}
		
		// Endereço
		if idx, exists := headerMap["endereco"]; exists && idx < len(row) {
			paciente.Endereco = strings.TrimSpace(row[idx])
		}
		
		// Município ID
		if idx, exists := headerMap["cod_municipio"]; exists && idx < len(row) {
			if municipioId, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				paciente.MunicipioID = municipioId
			}
		}
		
		// CEP
		if idx, exists := headerMap["cep"]; exists && idx < len(row) {
			paciente.CEP = strings.TrimSpace(row[idx])
		}
		
		// Telefone
		if idx, exists := headerMap["telefone"]; exists && idx < len(row) {
			paciente.Telefone = strings.TrimSpace(row[idx])
		}
		
		// Email
		if idx, exists := headerMap["email"]; exists && idx < len(row) {
			paciente.Email = strings.TrimSpace(row[idx])
		}
		
		// Contato de emergência
		if idx, exists := headerMap["contato_emergencia"]; exists && idx < len(row) {
			paciente.ContatoEmergencia = strings.TrimSpace(row[idx])
		}
		
		// Convênio
		if idx, exists := headerMap["convenio"]; exists && idx < len(row) {
			paciente.Convenio = strings.TrimSpace(row[idx])
		}
		
		// Número da carteira
		if idx, exists := headerMap["numero_carteira"]; exists && idx < len(row) {
			paciente.NumeroCarteira = strings.TrimSpace(row[idx])
		}
		
		// Status
		if idx, exists := headerMap["status"]; exists && idx < len(row) {
			paciente.Status = strings.TrimSpace(row[idx])
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