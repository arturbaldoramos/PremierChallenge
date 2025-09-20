package domain

import (
	"github.com/google/uuid"
)

type Medico struct {
	ID            int       `json:"id" db:"id"`
	UUID          uuid.UUID `json:"uuid" db:"uuid"`
	Nome          string    `json:"nome" db:"nome"`
	Especialidade string    `json:"especialidade" db:"especialidade"`
	CodMunicipio  string    `json:"cod_municipio" db:"cod_municipio"`
}
