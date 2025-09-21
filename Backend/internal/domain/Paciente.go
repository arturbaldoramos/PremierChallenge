package domain

import (
	"github.com/google/uuid"
)

type Paciente struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CPF          string    `json:"cpf" gorm:"uniqueIndex;size:11;not null"`
	Nome         string    `json:"nome" gorm:"not null;size:200"`
	Genero       string    `json:"genero" gorm:"size:1"`
	CodMunicipio string    `json:"cod_municipio" gorm:"size:20;index"`
	Bairro       string    `json:"bairro" gorm:"size:100"`
	Convenio     string    `json:"convenio" gorm:"size:10"`
	CID10        string    `json:"cid10" gorm:"size:10"`
}
