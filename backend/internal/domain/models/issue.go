package models

import (
	"main/pkg/utils"
	"slices"
	"time"
)

type IssueType string
type IssueId uint
type IssueIid string
type IssueExternalId string

const (
	IssueTypeEpic        IssueType = "EPIC"
	IssueTypeTask        IssueType = "TASK"
	IssueTypeIncident    IssueType = "INCIDENT"
	IssueTypeIssue       IssueType = "ISSUE"
	IssueTypeKeyResult   IssueType = "KEY_RESULT"
	IssueTypeObjective   IssueType = "OBJECTIVE"
	IssueTypeRequirement IssueType = "REQUIREMENT"
	IssueTypeTestCase    IssueType = "TEST_CASE"
	IssueTypeTicket      IssueType = "TICKET"
)

type Issue struct {
	Id           IssueId
	ExternalId   IssueExternalId
	Iid          IssueIid
	Title        string
	IssueType    IssueType
	Assignees    []User
	WebUrl       string
	Labels       []Label
	LabelHistory []LabelHistory
	Project      Project
	Release      *Release
	TaskType     *Label
	EstimateDev  *uint
	EstimateQA   *uint
}

func (i *Issue) Validate() error {
	return nil
}

func (i *Issue) SetHistory(history []LabelHistory) {
	i.LabelHistory = history
}

func (i *Issue) GetAddedHistory() []LabelHistory {
	labelIds := utils.Map(i.Labels, func(l Label) LabelId {
		return l.Id
	})

	var lastLabelIds = make([]LabelId, 0)
	if len(i.LabelHistory) > 0 {
		lastLabelIds = utils.Map(i.LabelHistory[len(i.LabelHistory)-1].Labels, func(id LabelId) LabelId {
			return id
		})
	}

	addedHistory := make([]LabelHistory, 0)

	if !slices.Equal(lastLabelIds, labelIds) {
		addedHistory = append(addedHistory, LabelHistory{
			Labels:    labelIds,
			CreatedAt: time.Now(),
		})
	}

	return addedHistory
}

func (i *Issue) BindToSprint(sprint *Sprint, assignee *User) (*IssueBinding, error) {
	priority := IssueBindingPriorityMedium
	binding := IssueBinding{
		Id:          0,
		Issue:       i,
		Sprint:      sprint,
		EstimateDev: nil,
		EstimateQA:  nil,
		BindStatus:  IssueBindStatusBacklog,
		Assignee:    assignee,
		Priority:    &priority,
	}

	if err := binding.Validate(); err != nil {
		return nil, err
	}

	return &binding, nil
}

func (i *Issue) GetAssignee() *User {
	if len(i.Assignees) > 0 {
		return &i.Assignees[0]
	}
	return nil
}
