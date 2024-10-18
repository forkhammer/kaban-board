package models

type Team struct {
	Id     int      `gorm:"id;primaryKey" json:"id"`
	Title  string   `gorm:"title;not null" json:"title"`
	Groups []*Group `gorm:"many2many:team_groups" json:"groups"`
}
