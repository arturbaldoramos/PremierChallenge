package domain

import (
	"github.com/google/uuid"
)

type Hospital struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID           uuid.UUID `json:"uuid" gorm:"type:uuid;uniqueIndex"`
	Nome           string    `json:"nome" gorm:"not null;size:200"`
	CEP            string    `json:"cep" gorm:"size:8;index"`
	Especialidades string    `json:"especialidades" gorm:"type:string;size:500"`
	LeitosTotais   int       `json:"leitos_totais" gorm:"default:0"`
	CodMunicipio   string    `json:"cod_municipio" gorm:"index;size:50"`
	Bairro         string    `json:"bairro" gorm:"size:100"`
}
