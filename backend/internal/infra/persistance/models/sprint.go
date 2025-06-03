package models

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	IsCompleted  bool      `gorm:"is_completed;default:false;not null"`
	Id           uint      `gorm:"id;primaryKey"`
	TeamId       uint      `gorm:"team_id;not null"`
	HoursPerUser uint      `gorm:"hours_per_user;default:0;not null"`
	Title        string    `gorm:"title;default:''"`
	StartDate    time.Time `gorm:"start_date;not null"`
	EndDate      time.Time `gorm:"end_date;not null"`
	Team         Team      `gorm:"foreignKey:team_id;not null"`
}
