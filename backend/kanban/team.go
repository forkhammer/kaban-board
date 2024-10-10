package kanban

type Team struct {
	Id     int      `gorm:"id;primaryKey" json:"id"`
	Title  string   `gorm:"title;not null" json:"title"`
	Groups []*Group `gorm:"many2many:team_groups" json:"groups"`
}

type UpdateTeamRequest struct {
	Title string `json:"title"`
}

type CreateTeamRequest struct {
	Title string `json:"title"`
}
