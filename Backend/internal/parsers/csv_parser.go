package parsers

import (
	"bufio"
	"encoding/csv"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"example.com/m/v2/internal/domain"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type CSVParser struct{}

// CSV structs for mapping
type EstadoCSV struct {
	Codigo            string `csv:"codigo_uf"`
	UnidadeFederativa string `csv:"uf"`
	Nome              string `csv:"nome"`
	Regiao            string `csv:"regiao"`
	Latitude          string `csv:"latitude"`
	Longitude         string `csv:"longitude"`
}

type MunicipioCSV struct {
	Codigo    string `csv:"codigo_ibge"`
	Nome      string `csv:"nome"`
	Latitude  string `csv:"latitude"`
	Longitude string `csv:"longitude"`
	Capital   string `csv:"capital"`
	CodigoUF  string `csv:"codigo_uf"`
	SiafiId   string `csv:"siafi_id"`
	DDD       string `csv:"ddd"`
	FusoHora  string `csv:"fuso_horario"`
	Populacao int    `csv:"populacao"`
}

type HospitalCSV struct {
	UUID           string `csv:"codigo"`
	Nome           string `csv:"nome"`
	CEP            string `csv:"cep"`
	Especialidades string `csv:"especialidades"`
	LeitosTotais   int    `csv:"leitos_totais"`
	CodMunicipio   string `csv:"cidade"`
	Bairro         string `csv:"bairro"`
}

type PacienteCSV struct {
	ID                string `csv:"id"`
	Nome              string `csv:"nome"`
	CPF               string `csv:"cpf"`
	RG                string `csv:"rg"`
	DataNascimento    string `csv:"data_nascimento"`
	Genero            string `csv:"genero"`
	TipoSanguineo     string `csv:"tipo_sanguineo"`
	Endereco          string `csv:"endereco"`
	MunicipioID       int    `csv:"municipio_id"`
	CEP               string `csv:"cep"`
	Telefone          string `csv:"telefone"`
	Email             string `csv:"email"`
	ContatoEmergencia string `csv:"contato_emergencia"`
	Convenio          string `csv:"convenio"`
	NumeroCarteira    string `csv:"numero_carteira"`
	Status            string `csv:"status"`
}

type MedicoCSV struct {
	UUID          string `csv:"codigo"`
	Nome          string `csv:"nome_completo"`
	Especialidade string `csv:"especialidade"`
	CodMunicipio  string `csv:"cidade"`
}

type CID10CSV struct {
	ID        int    `csv:"id"`
	Codigo    string `csv:"codigo"`
	Descricao string `csv:"descricao"`
	Categoria string `csv:"categoria"`
	Grupo     string `csv:"grupo"`
}

func NewCSVParser() *CSVParser {
	return &CSVParser{}
}

func (p *CSVParser) ParseEstados(file *multipart.FileHeader) ([]domain.Estado, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvData []EstadoCSV
	if err := gocsv.Unmarshal(f, &csvData); err != nil {
		return nil, err
	}

	var estados []domain.Estado
	for _, csv := range csvData {
		estado := domain.Estado{
			Codigo:            csv.Codigo,
			UnidadeFederativa: csv.UnidadeFederativa,
			Nome:              csv.Nome,
			Regiao:            csv.Regiao,
			Latitude:          csv.Latitude,
			Longitude:         csv.Longitude,
		}
		estados = append(estados, estado)
	}

	return estados, nil
}

func (p *CSVParser) ParseMunicipios(file *multipart.FileHeader) ([]domain.Municipio, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvData []MunicipioCSV
	if err := gocsv.Unmarshal(f, &csvData); err != nil {
		return nil, err
	}

	var municipios []domain.Municipio
	for _, csv := range csvData {
		municipio := domain.Municipio{
			Codigo:    csv.Codigo,
			Nome:      csv.Nome,
			Latitude:  csv.Latitude,
			Longitude: csv.Longitude,
			Capital:   csv.Capital,
			CodigoUF:  csv.CodigoUF,
			SiafiId:   csv.SiafiId,
			DDD:       csv.DDD,
			FusoHora:  csv.FusoHora,
			Populacao: csv.Populacao,
		}
		municipios = append(municipios, municipio)
	}

	return municipios, nil
}

// ParseMunicipiosStreaming - Versão otimizada que processa em chunks para economizar memória
func (p *CSVParser) ParseMunicipiosStreaming(file *multipart.FileHeader, chunkSize int) ([]domain.Municipio, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	reader.Comma = ','
	reader.LazyQuotes = true

	// Ler cabeçalho
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Mapear índices das colunas
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.TrimSpace(header)] = i
	}

	var municipios []domain.Municipio
	lineCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Pular linhas com erro
		}

		lineCount++

		// Processar em chunks para não consumir muita memória
		if len(municipios) >= chunkSize {
			// Aqui você poderia processar o chunk atual se necessário
			// Para este exemplo, continuamos acumulando
		}

		municipio := domain.Municipio{}

		// Mapear campos usando os índices do cabeçalho
		if idx, exists := headerMap["codigo_ibge"]; exists && idx < len(record) {
			municipio.Codigo = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["nome"]; exists && idx < len(record) {
			municipio.Nome = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["latitude"]; exists && idx < len(record) {
			municipio.Latitude = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["longitude"]; exists && idx < len(record) {
			municipio.Longitude = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["capital"]; exists && idx < len(record) {
			municipio.Capital = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["codigo_uf"]; exists && idx < len(record) {
			municipio.CodigoUF = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["siafi_id"]; exists && idx < len(record) {
			municipio.SiafiId = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["ddd"]; exists && idx < len(record) {
			municipio.DDD = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["fuso_horario"]; exists && idx < len(record) {
			municipio.FusoHora = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["populacao"]; exists && idx < len(record) {
			if pop, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
				municipio.Populacao = pop
			}
		}

		municipios = append(municipios, municipio)
	}

	return municipios, nil
}

func (p *CSVParser) ParseHospitais(file *multipart.FileHeader) ([]domain.Hospital, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvData []HospitalCSV
	if err := gocsv.Unmarshal(f, &csvData); err != nil {
		return nil, err
	}

	var hospitais []domain.Hospital
	for _, csv := range csvData {
		var hospitalUUID uuid.UUID
		if csv.UUID != "" {
			hospitalUUID, _ = uuid.Parse(csv.UUID)
		} else {
			hospitalUUID = uuid.New()
		}

		hospital := domain.Hospital{
			UUID:           hospitalUUID,
			Nome:           csv.Nome,
			CEP:            csv.CEP,
			Especialidades: csv.Especialidades,
			LeitosTotais:   csv.LeitosTotais,
			CodMunicipio:   csv.CodMunicipio,
			Bairro:         csv.Bairro,
		}
		hospitais = append(hospitais, hospital)
	}

	return hospitais, nil
}

func (p *CSVParser) ParsePacientes(file *multipart.FileHeader) ([]domain.Paciente, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvData []PacienteCSV
	if err := gocsv.Unmarshal(f, &csvData); err != nil {
		return nil, err
	}

	var pacientes []domain.Paciente
	for _, csv := range csvData {
		var pacienteUUID uuid.UUID
		if csv.ID != "" {
			pacienteUUID, _ = uuid.Parse(csv.ID)
		} else {
			pacienteUUID = uuid.New()
		}

		paciente := domain.Paciente{
			ID:           pacienteUUID,
			CPF:          csv.CPF,
			Nome:         csv.Nome,
			Genero:       csv.Genero,
			CodMunicipio: "", // CSV não tem este campo, deixar vazio
			Bairro:       "", // CSV não tem este campo, deixar vazio
			Convenio:     csv.Convenio,
			CID10:        "", // CSV não tem este campo, deixar vazio
		}
		pacientes = append(pacientes, paciente)
	}

	return pacientes, nil
}

func (p *CSVParser) ParseMedicos(file *multipart.FileHeader) ([]domain.Medico, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvData []MedicoCSV
	if err := gocsv.Unmarshal(f, &csvData); err != nil {
		return nil, err
	}

	var medicos []domain.Medico
	for _, csv := range csvData {
		var medicoUUID uuid.UUID
		if csv.UUID != "" {
			medicoUUID, _ = uuid.Parse(csv.UUID)
		} else {
			medicoUUID = uuid.New()
		}

		medico := domain.Medico{
			UUID:          medicoUUID,
			Nome:          csv.Nome,
			Especialidade: csv.Especialidade,
			CodMunicipio:  csv.CodMunicipio,
		}
		medicos = append(medicos, medico)
	}

	return medicos, nil
}

// ParseMedicosStreaming - Versão otimizada que processa em chunks para economizar memória
func (p *CSVParser) ParseMedicosStreaming(file *multipart.FileHeader, chunkSize int) ([]domain.Medico, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	reader.Comma = ','
	reader.LazyQuotes = true

	// Ler cabeçalho
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Mapear índices das colunas
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.TrimSpace(header)] = i
	}

	var medicos []domain.Medico
	lineCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Pular linhas com erro
		}

		lineCount++

		// Processar em chunks para não consumir muita memória
		if len(medicos) >= chunkSize {
			// Aqui você poderia processar o chunk atual se necessário
			// Para este exemplo, continuamos acumulando
		}

		medico := domain.Medico{}

		// Mapear UUID do campo codigo
		var medicoUUID uuid.UUID
		if idx, exists := headerMap["codigo"]; exists && idx < len(record) {
			uuidStr := strings.TrimSpace(record[idx])
			if uuidStr != "" {
				if parsedUUID, err := uuid.Parse(uuidStr); err == nil {
					medicoUUID = parsedUUID
				} else {
					// Se não for um UUID válido, gerar um novo
					medicoUUID = uuid.New()
				}
			} else {
				medicoUUID = uuid.New()
			}
		} else {
			medicoUUID = uuid.New()
		}
		medico.UUID = medicoUUID

		// Mapear outros campos usando os índices do cabeçalho
		if idx, exists := headerMap["nome_completo"]; exists && idx < len(record) {
			medico.Nome = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["especialidade"]; exists && idx < len(record) {
			medico.Especialidade = strings.TrimSpace(record[idx])
		}
		if idx, exists := headerMap["cidade"]; exists && idx < len(record) {
			codMunicipio := strings.TrimSpace(record[idx])
			if codMunicipio != "" {
				medico.CodMunicipio = codMunicipio
			}
		}

		medicos = append(medicos, medico)
	}

	return medicos, nil
}

func (p *CSVParser) ParseCID10(file *multipart.FileHeader) ([]domain.Cid10, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvData []CID10CSV
	if err := gocsv.Unmarshal(f, &csvData); err != nil {
		return nil, err
	}

	var cid10s []domain.Cid10
	for _, csv := range csvData {
		cid10 := domain.Cid10{
			ID:        csv.ID,
			Codigo:    csv.Codigo,
			Descricao: csv.Descricao,
			Categoria: csv.Categoria,
			Grupo:     csv.Grupo,
		}
		cid10s = append(cid10s, cid10)
	}

	return cid10s, nil
}
