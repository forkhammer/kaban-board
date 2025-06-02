package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LabelHistoryProcess string

type LabelHistory struct {
	gorm.Model
	IssueId uint                        `gorm:"issue_id;not null"`
	Issue   Issue                       `gorm:"foreignKey:IssueId;not null"`
	Labels  datatypes.JSONSlice[string] `gorm:"labels;not null"`
}
