package models

import "gorm.io/gorm"

type Issue struct {
	gorm.Model
	Id          uint     `gorm:"primarykey"`
	Iid         string   `gorm:"iid"`
	Title       string   `gorm:"title;not null"`
	IssueType   string   `gorm:"issue_type;not null"`
	Assignees   []User   `gorm:"many2many:assignees;"`
	WebUrl      string   `gorm:"web_url;not null"`
	Labels      []Label  `gorm:"many2many:issue_labels;"`
	ProjectId   uint     `gorm:"project_id;not null"`
	Project     Project  `gorm:"foreignKey:ProjectId;not null"`
	ReleaseId   *string  `gorm:"release_id"`
	Release     *Release `gorm:"foreignKey:ReleaseId"`
	TaskTypeId  *string  `gorm:"task_type_id"`
	TaskType    *Label   `gorm:"foreignKey:TaskTypeId"`
	EstimateDev *uint    `gorm:"estimate_dev"`
	EstimateQA  *uint    `gorm:"estimate_qa"`
}
