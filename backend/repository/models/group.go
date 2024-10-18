package models

type Group struct {
	Id   uint   `gorm:"id;primaryKey" json:"id"`
	Name string `gorm:"name;not null" json:"title"`
}
