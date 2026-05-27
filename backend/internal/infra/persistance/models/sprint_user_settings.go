package models

import "time"

type SprintUserSettings struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Id           uint `gorm:"id;primaryKey"`
	SprintId     uint `gorm:"sprint_id;not null;uniqueIndex:sprint_user_idx"`
	UserId       uint `gorm:"user_id;not null;uniqueIndex:sprint_user_idx"`
	User         User `gorm:"foreignKey:UserId"`
	HoursPerUser uint `gorm:"hours_per_user;default:0;not null"`
}
