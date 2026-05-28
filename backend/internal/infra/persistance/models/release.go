package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Release struct {
	gorm.Model
	Id        uint                    `gorm:"primaryKey"`
	Iid       string                  `gorm:"iid;not null"`
	Title     string                  `gorm:"title;not null"`
	ProjectId uint                    `gorm:"project_id;not null"`
	Project   Project                 `gorm:"foreignKey:ProjectId;not null"`
	WebPath   string                  `gorm:"web_path;not null;default:''"`
	Index     datatypes.JSONSlice[int64] `gorm:"column:sort_index"`
}
