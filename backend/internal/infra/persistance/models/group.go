package models

type Group struct {
	Id   uint   `gorm:"id;primaryKey"`
	Name string `gorm:"name;not null"`
}
