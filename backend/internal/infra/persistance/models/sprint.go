package models

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	Id           uint      `gorm:"id;primaryKey"`
	Status       string    `gorm:"status;default:waiting;not null"`
	TeamId       uint      `gorm:"team_id;not null"`
	HoursPerUser uint      `gorm:"hours_per_user;default:0;not null"`
	Title        string    `gorm:"title;default:''"`
	StartDate    time.Time `gorm:"start_date;not null"`
	EndDate      time.Time `gorm:"end_date;not null"`
	Team         Team      `gorm:"foreignKey:team_id;not null"`
}
