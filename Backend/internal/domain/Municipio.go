package domain

type Municipio struct {
	ID        int    `json:"id" db:"id"`
	Codigo    string `json:"codigo" db:"codigo"` // Código IBGE
	Nome      string `json:"nome" db:"nome"`
	Latitude  string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
	Capital   string `json:"capital" db:"capital"`
	CodigoUF  string `json:"codigo_uf" db:"codigo_uf"`
	SiafiId   string `json:"siafi_id" db:"siafi_id"`
	DDD       string `json:"ddd" db:"ddd"`
	FusoHora  string `json:"fuso_hora" db:"fuso_hora"`
	Populacao int    `json:"populacao" db:"populacao"`
}
