package models

import (
	"time"

	"gorm.io/gorm"
)

type IssueBindingHistory struct {
	gorm.Model
	IssueBindingId uint         `gorm:"issue_binding_id;not null;index"`
	IssueBinding   IssueBinding `gorm:"foreignKey:IssueBindingId;not null"`
	IssueId        uint         `gorm:"issue_id;not null;index"`
	Issue          Issue        `gorm:"foreignKey:IssueId;not null"`
	EstimateDev    *uint        `gorm:"estimate_dev"`
	EstimateQA     *uint        `gorm:"estimate_qa"`
	BindStatus     string       `gorm:"bind_status;not null"`
	Priority       *string      `gorm:"priority"`
	AssigneeId     *uint        `gorm:"assignee_id"`
	Assignee       *User        `gorm:"foreignKey:AssigneeId"`
	DeletedAt      *time.Time   `gorm:"deleted_at;index"`
	CreatedAt      time.Time    `gorm:"created_at;autoCreateTime"`
}
