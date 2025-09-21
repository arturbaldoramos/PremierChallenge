package domain

type Municipio struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Codigo    string `json:"codigo" gorm:"uniqueIndex;not null;size:10"`
	Nome      string `json:"nome" gorm:"not null;size:100"`
	Latitude  string `json:"latitude" gorm:"size:20"`
	Longitude string `json:"longitude" gorm:"size:20"`
	Capital   string `json:"capital" gorm:"size:1"`
	CodigoUF  string `json:"codigo_uf" gorm:"index;size:10"`
	SiafiId   string `json:"siafi_id" gorm:"size:10"`
	DDD       string `json:"ddd" gorm:"size:3"`
	FusoHora  string `json:"fuso_hora" gorm:"size:50"`
	Populacao int    `json:"populacao" gorm:"default:0"`
}
