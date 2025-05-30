package models

type Team struct {
	Id     uint     `gorm:"id;primaryKey"`
	Title  string   `gorm:"title;not null"`
	Groups []*Group `gorm:"many2many:team_groups"`
}
