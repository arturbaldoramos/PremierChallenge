package domain

type Cid10 struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Codigo    string `json:"codigo" gorm:"uniqueIndex;not null;size:10"`
	Descricao string `json:"descricao" gorm:"not null;type:text"`
	Categoria string `json:"categoria" gorm:"index;size:100"`
	Grupo     string `json:"grupo" gorm:"index;size:100"`
}
