package kanban

type Label struct {
	Id        string  `gorm:"id;primaryKey" json:"id"`
	Name      string  `gorm:"name;not null" json:"name"`
	Color     string  `gorm:"color" json:"color"`
	TextColor string  `gorm:"text_color" json:"textColor"`
	AltName   *string `gorm:"alt_name" json:"altName"`
}

type KanbanLabel struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	AltName *string `json:"altName"`
}

type UpdateKanbanLabelRequest struct {
	AltName *string `json:"altName"`
}
