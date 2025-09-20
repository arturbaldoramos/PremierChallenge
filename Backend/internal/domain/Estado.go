package domain

type Estado struct {
	ID                int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Codigo            string `json:"codigo" gorm:"uniqueIndex;not null"`
	UnidadeFederativa string `json:"unidade_federativa" gorm:"not null"`
	Nome              string `json:"nome" gorm:"not null"`
	Regiao            string `json:"regiao" gorm:"not null"`
	Latitude          string `json:"latitude" gorm:"not null"`
	Longitude         string `json:"longitude" gorm:"not null"`
}
