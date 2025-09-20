package domain

type Estado struct {
	ID                int    `json:"id" db:"id"`
	Codigo            string `json:"codigo" db:"codigo"`                         // Cod IBGE
	UnidadeFederativa string `json:"unidade_federativa" db:"unidade_federativa"` // UF (SP, RJ, etc.)
	Nome              string `json:"nome" db:"nome"`
	Regiao            string `json:"regiao" db:"regiao"`
	Latitude          string `json:"latitude" db:"latitude"`
	Longitude         string `json:"longitude" db:"longitude"`
}
