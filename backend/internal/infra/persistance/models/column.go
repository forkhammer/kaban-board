package models

import (
	"gorm.io/datatypes"
)

type Column struct {
	Id     int                         `gorm:"id;primaryKey"`
	Name   string                      `gorm:"name"`
	Labels datatypes.JSONSlice[string] `gorm:"labels"`
	TeamId *uint                       `gorm:"team_id"`
	Team   *Team                       `gorm:"foreignKey:team_id"`
	Order  *int                        `gorm:"order;not null;default:10"`
}
