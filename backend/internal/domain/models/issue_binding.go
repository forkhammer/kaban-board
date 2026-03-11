package models

import "github.com/go-playground/validator/v10"

type IssueBindingId uint

type IssueBindingStatus string

const (
	IssueBindStatusBacklog    IssueBindingStatus = "backlog"
	IssueBindStatusInProgress IssueBindingStatus = "in_progress"
	IssueBindStatusDone       IssueBindingStatus = "done"
)

type IssueBindingPriority string

const (
	IssueBindingPriorityLowest   IssueBindingPriority = "lowest"
	IssueBindingPriorityLow      IssueBindingPriority = "low"
	IssueBindingPriorityMedium   IssueBindingPriority = "medium"
	IssueBindingPriorityHigh     IssueBindingPriority = "high"
	IssueBindingPriorityCritical IssueBindingPriority = "critical"
)

type IssueBinding struct {
	Id          IssueBindingId
	Sprint      *Sprint `validate:"required"`
	Issue       *Issue  `validate:"required"`
	EstimateDev *uint
	EstimateQA  *uint
	BindStatus  IssueBindingStatus
	Priority    *IssueBindingPriority `validate:"omitnil,oneof=lowest low medium high critical"`
	Assignee    *User
	Comment     *string
	Release     *Release
	Epic        *Epic
	Order       string
}

func (ib *IssueBinding) Validate() error {
	validator := validator.New()
	if err := validator.Struct(ib); err != nil {
		return err
	}

	return nil
}

func (ib *IssueBinding) CanUpdate(account *Account) bool {
	if account == nil {
		return false
	}

	if account.Role == AccountRoleAdmin {
		return true
	}

	if account.Role == AccountRoleEmployee && account.GitlabID != nil && ib.Assignee != nil {
		if ib.Assignee.Id == *account.GitlabID {
			return true
		}
	}

	return false
}

func (ib *IssueBinding) CanManage(account *Account) bool {
	if account == nil {
		return false
	}

	if account.Role == AccountRoleAdmin {
		return true
	}

	return false
}
