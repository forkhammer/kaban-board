package kanban

type Team struct {
	Id    int    `gorm:"id;primaryKey" json:"id"`
	Title string `gorm:"title;not null" json:"title"`
}

type UpdateTeamRequest struct {
	Title string `json:"title"`
}

type CreateTeamRequest struct {
	Title string `json:"title"`
}
