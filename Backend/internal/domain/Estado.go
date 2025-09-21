package domain

type Estado struct {
	ID                int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Codigo            string `json:"codigo" gorm:"uniqueIndex;not null;size:2"`
	UnidadeFederativa string `json:"unidade_federativa" gorm:"not null;size:2"`
	Nome              string `json:"nome" gorm:"not null;size:100"`
	Regiao            string `json:"regiao" gorm:"not null;size:50"`
	Latitude          string `json:"latitude" gorm:"size:20"`
	Longitude         string `json:"longitude" gorm:"size:20"`
}
