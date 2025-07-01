package models

import (
	"fmt"
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

func (i *Issue) UnbindFromSprint(bindingId IssueBindingId) error {
	bindingCount := len(i.SprintBindings)
	if bindingCount == 0 {
		return fmt.Errorf("Issue %d has no binding %d", i.Id, bindingId)
	}
	i.SprintBindings = utils.Filter(i.SprintBindings, func(binding IssueBinding) bool {
		return binding.Id != bindingId
	})
	if len(i.SprintBindings) == bindingCount {
		return fmt.Errorf("Issue %d has no binding %d", i.Id, bindingId)
	}
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

func (i *Issue) SetEstimateDev(estimate *uint) error {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.EstimateDev = estimate
			return i.Validate()
		}
	}

	i.EstimateDev = estimate
	return i.Validate()
}

func (i *Issue) SetEstimateQA(estimate *uint) error {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.EstimateQA = estimate
			return i.Validate()
		}
	}

	i.EstimateQA = estimate
	return i.Validate()
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

func (i *Issue) SetBindStatus(status IssueBindingStatus) error {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.BindStatus = status
			return i.Validate()
		}
	}
	return nil
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

func (i *Issue) GetComment() *string {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			return binding.Comment
		}
	}
	return nil
}

func (i *Issue) SetComment(comment *string) error {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.Comment = comment
			return i.Validate()
		}
	}
	return nil
}

func (i *Issue) GetPriority() *IssueBindingPriority {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			return binding.Priority
		}
	}
	return nil
}

func (i *Issue) SetPriority(priority *IssueBindingPriority) error {
	if i.contextBindingId != nil {
		binding := i.getBindingById(*i.contextBindingId)
		if binding != nil {
			binding.Priority = priority
			return i.Validate()
		}
	}
	return nil
}

func (i *Issue) getBindingById(id IssueBindingId) *IssueBinding {
	for index, binding := range i.SprintBindings {
		if binding.Id == id {
			return &i.SprintBindings[index]
		}
	}
	return nil
}
