package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type IssueBindingHistoryId uint

type IssueBindingHistory struct {
	Id             IssueBindingHistoryId
	IssueBindingId IssueBindingId
	IssueId        IssueId
	SprintId       SprintId
	EstimateDev    *uint
	EstimateQA     *uint
	BindStatus     IssueBindingStatus
	Priority       *IssueBindingPriority `validate:"omitnil,oneof=lowest low medium high critical"`
	Assignee       *User
	RemovedAt      *time.Time
	CreatedAt      time.Time
}

func (ibh *IssueBindingHistory) Validate() error {
	validator := validator.New()
	if err := validator.Struct(ibh); err != nil {
		return err
	}

	return nil
}
