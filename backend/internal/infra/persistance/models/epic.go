package models

import "gorm.io/gorm"

type Epic struct {
	gorm.Model
	Id         uint    `gorm:"primaryKey"`
	ExternalId string  `gorm:"column:external_id;uniqueIndex"`
	Iid        string  `gorm:"iid"`
	Title      string  `gorm:"title;not null"`
	ProjectId  uint    `gorm:"project_id;not null"`
	Project    Project `gorm:"foreignKey:ProjectId;not null"`
}
