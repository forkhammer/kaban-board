package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Project struct {
	IsVisible bool `gorm:"is_visible;default:true;not null"`
	gorm.Model
	Id     uint                       `gorm:"primaryKey"`
	Name   string                     `gorm:"name;not null"`
	TeamId *int                       `gorm:"team_id"`
	Team   Team                       `gorm:"foreignKey:team_id"`
	Users  datatypes.JSONSlice[int64] `gorm:"users"`
}
