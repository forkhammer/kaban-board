package kanban

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

type UpdateColumnRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	TeamId *int     `json:"team_id"`
}

type CreateColumnRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	TeamId *int     `json:"team_id"`
}

type SetColumnOrderRequest struct {
	Id    int `json:"id"`
	Order int `json:"order"`
}
