package kanban

import (
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Id        uint          `gorm:"primaryKey" json:"id"`
	Name      string        `json:"name"`
	IsVisible bool          `gorm:"is_visible;default:true;not null" json:"is_visible"`
	TeamId    *int          `gorm:"team_id" json:"team_id"`
	Team      Team          `gorm:"foreignKey:team_id" json:"-"`
	Users     pq.Int64Array `gorm:"users;type:integer[]" json:"-"`
}

type SetTeamRequest struct {
	TeamId *int `json:"team_id"`
}
