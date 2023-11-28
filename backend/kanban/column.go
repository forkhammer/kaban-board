package kanban

import (
	"github.com/lib/pq"
)

type Column struct {
	Id     int            `gorm:"id;primaryKey" json:"id"`
	Name   string         `gorm:"name" json:"name"`
	Labels pq.StringArray `gorm:"labels; type:varchar[]" json:"labels"`
	TeamId *int           `gorm:"team_id" json:"team_id"`
	Team   *Team          `gorm:"foreignKey:team_id" json:"team"`
	Order  *int           `gorm:"order;not null;default:10" json:"order"`
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
