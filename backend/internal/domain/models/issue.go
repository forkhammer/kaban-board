package models

import (
	"main/pkg/utils"
	"slices"
	"time"
)

type IssueType string
type IssueId uint
type IssueIid string

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
	Id               IssueId
	Iid              IssueIid
	Title            string
	IssueType        IssueType
	Assignees        []User
	WebUrl           string
	Labels           []Label
	LabelHistory     []LabelHistory
	Project          Project
	Release          *Release
	TaskType         *Label
	EstimateDev      *uint
	EstimateQA       *uint
	SprintBindings   []IssueBinding
	contextSprintId  *SprintId
	contextBindingId *IssueBindingId
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

func (i *Issue) BindToSprint(sprint *Sprint) error {
	binding := IssueBinding{
		Id:          0,
		Issue:       i,
		Sprint:      sprint,
		EstimateDev: nil,
		EstimateQA:  nil,
		BindStatus:  IssueBindStatusBacklog,
		Assignee:    nil,
	}

	if err := binding.Validate(); err != nil {
		return err
	}

	i.SprintBindings = append(i.SprintBindings, binding)
	return nil
}

func (i *Issue) SetContext(bindingId *IssueBindingId) {
	i.contextBindingId = bindingId
}

func (i *Issue) GetContextBindingId() *IssueBindingId {
	return i.contextBindingId
}

func (i *Issue) GetEstimateDev() *uint {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			return binding.EstimateDev
		}
	}

	return i.EstimateDev
}

func (i *Issue) GetEstimateQA() *uint {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			return binding.EstimateQA
		}
	}

	return i.EstimateQA
}

func (i *Issue) SetEstimateDev(estimate *uint) {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.EstimateDev = estimate
			return
		}
	}

	i.EstimateDev = estimate
}

func (i *Issue) SetEstimateQA(estimate *uint) {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.EstimateQA = estimate
			return
		}
	}

	i.EstimateQA = estimate
}

func (i *Issue) GetBindStatus() *IssueBindingStatus {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			return &binding.BindStatus
		}
	}
	return nil
}

func (i *Issue) SetBindStatus(status IssueBindingStatus) {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.BindStatus = status
			return
		}
	}
}

func (i *Issue) GetAssignee() *User {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil && binding.Assignee != nil {
			return binding.Assignee
		}
	}

	if len(i.Assignees) > 0 {
		return &i.Assignees[0]
	}

	return nil
}

func (i *Issue) SetAssignee(assignee *User) {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.Assignee = assignee
			return
		}
	}
}

func (i *Issue) getBindingById(id IssueBindingId) *IssueBinding {
	for index, binding := range i.SprintBindings {
		if binding.Id == id {
			return &i.SprintBindings[index]
		}
	}
	return nil
}
