package models

type IssueBinding struct {
	Id          uint    `gorm:"id;primaryKey"`
	SprintId    uint    `gorm:"sprint_id;not null"`
	Sprint      *Sprint `gorm:"foreignKey:sprint_id;not null"`
	IssueId     uint    `gorm:"issue_id;not null"`
	Issue       *Issue  `gorm:"foreignKey:issue_id;not null"`
	EstimateDev *uint   `gorm:"estimate_dev"`
	EstimateQA  *uint   `gorm:"estimate_qa"`
	BindStatus  string  `gorm:"bind_status;not null"`
	Priority    *string `gorm:"priority"`
	AssigneeId  *uint   `gorm:"assignee_id"`
	Assignee    *User   `gorm:"foreignKey:assignee_id"`
	Comment     *string `gorm:"comment"`
}
