package parsers

import (
	"mime/multipart"
	"time"

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
	ID             int    `csv:"id"`
	UUID           string `csv:"uuid"`
	Nome           string `csv:"nome"`
	CEP            string `csv:"cep"`
	Especialidades string `csv:"especialidades"`
	LeitosTotais   int    `csv:"leitos_totais"`
	CodMunicipio   string `csv:"cod_municipio"`
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
	ID            string `csv:"id"`
	Nome          string `csv:"nome"`
	Especialidade string `csv:"especialidade"`
	CodMunicipio  string `csv:"cod_municipio"`
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
			ID:           csv.ID,
			UUID:         hospitalUUID,
			Nome:         csv.Nome,
			CEP:          csv.CEP,
			LeitosTotais: csv.LeitosTotais,
			CodMunicipio: csv.CodMunicipio,
			Bairro:       csv.Bairro,
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

		dataNascimento, _ := time.Parse("2006-01-02", csv.DataNascimento)

		paciente := domain.Paciente{
			ID:                pacienteUUID,
			Nome:              csv.Nome,
			CPF:               csv.CPF,
			RG:                csv.RG,
			DataNascimento:    dataNascimento,
			Genero:            csv.Genero,
			TipoSanguineo:     csv.TipoSanguineo,
			Endereco:          csv.Endereco,
			MunicipioID:       csv.MunicipioID,
			CEP:               csv.CEP,
			Telefone:          csv.Telefone,
			Email:             csv.Email,
			ContatoEmergencia: csv.ContatoEmergencia,
			Convenio:          csv.Convenio,
			NumeroCarteira:    csv.NumeroCarteira,
			Status:            csv.Status,
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
		if csv.ID != "" {
			medicoUUID, _ = uuid.Parse(csv.ID)
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
