package domain

import (
	"github.com/google/uuid"
)

type Medico struct {
	UUID          uuid.UUID `json:"uuid" gorm:"type:uuid;uniqueIndex"`
	Nome          string    `json:"nome" gorm:"not null;size:200"`
	Especialidade string    `json"especialidade" gorm:"index;size:100"`
	CodMunicipio  string    `json:"cod_municipio" gorm:"size:50"`
}
