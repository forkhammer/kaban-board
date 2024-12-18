package models

import (
	"gorm.io/datatypes"
)

type Column struct {
	Id     int                         `gorm:"id;primaryKey" json:"id"`
	Name   string                      `gorm:"name" json:"name"`
	Labels datatypes.JSONSlice[string] `gorm:"labels" json:"labels"`
	TeamId *int                        `gorm:"team_id" json:"team_id"`
	Team   *Team                       `gorm:"foreignKey:team_id" json:"team"`
	Order  *int                        `gorm:"order;not null;default:10" json:"order"`
}
