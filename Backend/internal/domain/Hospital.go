package domain

import "github.com/google/uuid"

type Hospital struct {
	ID             int       `json:"id" db:"id"`
	UUID           uuid.UUID `json:"uuid" db:"uuid"`
	Nome           string    `json:"nome" db:"nome"`
	CEP            string    `json:"cep" db:"cep"`
	Especialidades []string  `json:"especialidades" db:"especialidades"` // JSONB
	LeitosTotais   int       `json:"leitos_totais" db:"leitos_totais"`
	CodMunicipio   string    `json:"cod_municipio" db:"cod_municipio"`
	Bairro         string    `json:"bairro" db:"bairro"`
}
