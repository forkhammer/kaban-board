package models

import "github.com/go-playground/validator/v10"

type IssueBindingId uint

type IssueBindingStatus string

const (
	IssueBindStatusBacklog    IssueBindingStatus = "backlog"
	IssueBindStatusInProgress IssueBindingStatus = "in_progress"
	IssueBindStatusDone       IssueBindingStatus = "done"
)

type IssueBinding struct {
	Id          IssueBindingId
	Sprint      *Sprint `validate:"required"`
	Issue       *Issue  `validate:"required"`
	EstimateDev *uint
	EstimateQA  *uint
	BindStatus  IssueBindingStatus
	Assignee    *User
	Comment     *string
}

func (ib *IssueBinding) Validate() error {
	validator := validator.New()
	if err := validator.Struct(ib); err != nil {
		return err
	}

	return nil
}
