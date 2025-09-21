package domain

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Hospital struct {
	ID                     int            `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID                   uuid.UUID      `json:"uuid" gorm:"type:uuid;uniqueIndex"`
	Nome                   string         `json:"nome" gorm:"not null;size:200"`
	CEP                    string         `json:"cep" gorm:"size:8;index"`
	Especialidades         datatypes.JSON `json:"especialidades" gorm:"type:jsonb"`
	EspecialidadePrimaria  string         `json:"especialidade_primaria" gorm:"size:100"`
	EspecialidadeSecundaria string        `json:"especialidade_secundaria" gorm:"size:100"`
	EspecialidadeTerciaria string         `json:"especialidade_terciaria" gorm:"size:100"`
	EspecialidadeQuaternaria string       `json:"especialidade_quaternaria" gorm:"size:100"`
	LeitosTotais           int            `json:"leitos_totais" gorm:"default:0"`
	CodMunicipio           string         `json:"cod_municipio" gorm:"index;size:50"`
	Bairro                 string         `json:"bairro" gorm:"size:100"`
}
