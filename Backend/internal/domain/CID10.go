package domain

type CID10 struct {
	ID        int    `json:"id" db:"id"`
	Codigo    string `json:"codigo" db:"codigo"` // A00.0, B15.9, etc.
	Descricao string `json:"descricao" db:"descricao"`
	Categoria string `json:"categoria" db:"categoria"`
	Grupo     string `json:"grupo" db:"grupo"`
}
