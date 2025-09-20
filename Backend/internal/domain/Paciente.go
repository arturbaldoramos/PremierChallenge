package domain

import (
	"time"

	"github.com/google/uuid"
)

type Paciente struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Nome              string    `json:"nome" db:"nome"`
	CPF               string    `json:"cpf" db:"cpf"`
	RG                string    `json:"rg" db:"rg"`
	DataNascimento    time.Time `json:"data_nascimento" db:"data_nascimento"`
	Genero            string    `json:"genero" db:"genero"`
	TipoSanguineo     string    `json:"tipo_sanguineo" db:"tipo_sanguineo"`
	Endereco          string    `json:"endereco" db:"endereco"`
	MunicipioID       int       `json:"municipio_id" db:"municipio_id"`
	CEP               string    `json:"cep" db:"cep"`
	Telefone          string    `json:"telefone" db:"telefone"`
	Email             string    `json:"email" db:"email"`
	ContatoEmergencia string    `json:"contato_emergencia" db:"contato_emergencia"`
	Convenio          string    `json:"convenio" db:"convenio"`
	NumeroCarteira    string    `json:"numero_carteira" db:"numero_carteira"`
	Status            string    `json:"status" db:"status"` // ativo, inativo
}
