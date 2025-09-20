package domain

import (
	"time"

	"github.com/google/uuid"
)

type Paciente struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Nome              string    `json:"nome" gorm:"not null;size:200"`
	CPF               string    `json:"cpf" gorm:"uniqueIndex;size:11"`
	RG                string    `json:"rg" gorm:"index;size:20"`
	DataNascimento    time.Time `json:"data_nascimento" gorm:"not null"`
	Genero            string    `json:"genero" gorm:"size:1"`
	TipoSanguineo     string    `json:"tipo_sanguineo" gorm:"size:3"`
	Endereco          string    `json:"endereco" gorm:"size:300"`
	MunicipioID       int       `json:"municipio_id" gorm:"index"`
	CEP               string    `json:"cep" gorm:"size:8"`
	Telefone          string    `json:"telefone" gorm:"size:20"`
	Email             string    `json:"email" gorm:"index;size:100"`
	ContatoEmergencia string    `json:"contato_emergencia" gorm:"size:100"`
	Convenio          string    `json:"convenio" gorm:"size:100"`
	NumeroCarteira    string    `json:"numero_carteira" gorm:"size:50"`
	Status            string    `json:"status" gorm:"default:ativo;size:20"`
}
