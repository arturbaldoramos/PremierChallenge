package domain

type Cid10 struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Codigo    string `json:"codigo" gorm:"uniqueIndex;not null;size:255"`
	Descricao string `json:"descricao" gorm:"type:text"`
	Categoria string `json:"categoria" gorm:"index;size:255"`
	Grupo     string `json:"grupo" gorm:"index;size:255"`
}

// TableName especifica o nome da tabela para o GORM
func (Cid10) TableName() string {
	return "cid10"
}
