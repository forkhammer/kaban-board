package dto

import (
	"fmt"
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type IssueBindingDto struct {
	Id          string      `json:"id"`
	Iid         string      `json:"iid"`
	Title       string      `json:"title"`
	IssueType   string      `json:"type"`
	Assignees   []UserDto   `json:"assignees"`
	Assignee    *UserDto    `json:"assignee"`
	WebUrl      string      `json:"webUrl"`
	Labels      []LabelDto  `json:"labels"`
	ProjectId   int         `json:"projectId"`
	ProjectName *string     `json:"projectName"`
	Release     *ReleaseDto `json:"release"`
	Epic        *EpicDto    `json:"epic"`
	TaskType    *LabelDto   `json:"taskType"`
	EstimateDev *uint       `json:"estimateDev"`
	EstimateQA  *uint       `json:"estimateQA"`
	BindingId   uint        `json:"bindingId"`
	BindStatus  string      `json:"bindStatus"`
	Priority    *string     `json:"priority"`
	Comment     *string     `json:"comment"`
	CanUpdate   bool        `json:"can_update"`
}

type IssueBindingPageDto struct {
	Count      int               `json:"count"`
	Pages      int               `json:"pages"`
	Page       int               `json:"page"`
	StartIndex int               `json:"start_index"`
	EndIndex   int               `json:"end_index"`
	Results    []IssueBindingDto `json:"results"`
}

func SerializeIssueBindingPage(issueBindingPage *usecases.IssueBindingPage, account *domain.Account) *IssueBindingPageDto {
	totalPages := int(issueBindingPage.Count / issueBindingPage.Limit)
	if issueBindingPage.Count%issueBindingPage.Limit > 0 {
		totalPages++
	}

	startIndex := int((issueBindingPage.Page - 1) * issueBindingPage.Limit)
	endIndex := startIndex + len(issueBindingPage.Results)

	return &IssueBindingPageDto{
		Count:      int(issueBindingPage.Count),
		Pages:      totalPages,
		Page:       int(issueBindingPage.Page),
		StartIndex: startIndex,
		EndIndex:   endIndex,
		Results:    SerializeIssueBindings(issueBindingPage.Results, account),
	}
}

func SerializeIssueBindings(bindings []domain.IssueBinding, account *domain.Account) []IssueBindingDto {
	return utils.Map(bindings, func(binding domain.IssueBinding) IssueBindingDto {
		return *SerializeIssueBinding(&binding, account)
	})
}

func SerializeIssueBinding(binding *domain.IssueBinding, account *domain.Account) *IssueBindingDto {
	dto := &IssueBindingDto{
		Id:          fmt.Sprintf("%d", binding.Issue.Id),
		EstimateDev: binding.EstimateDev,
		EstimateQA:  binding.EstimateQA,
		BindingId:   (uint)(binding.Id),
		BindStatus:  (string)(binding.BindStatus),
		CanUpdate:   binding.CanUpdate(account),
	}

	if binding.Priority != nil {
		priority := string(*binding.Priority)
		dto.Priority = &priority
	}

	dto.Comment = binding.Comment

	if binding.Assignee != nil {
		dto.Assignee = SerializeUser(binding.Assignee)
	}

	if binding.Release != nil {
		dto.Release = SerializeRelease(binding.Release)
	}

	if binding.Epic != nil {
		dto.Epic = SerializeEpic(binding.Epic)
	}

	if binding.Issue != nil {
		issue := binding.Issue
		dto.Iid = string(issue.Iid)
		dto.Title = issue.Title
		dto.IssueType = string(issue.IssueType)
		dto.Assignees = SerializeUsers(issue.Assignees)
		dto.WebUrl = issue.WebUrl
		dto.Labels = SerializeLabels(issue.Labels)
		dto.ProjectId = int(issue.Project.Id)
		dto.ProjectName = &issue.Project.Name
		dto.TaskType = func() *LabelDto {
			if issue.TaskType == nil {
				return nil
			}
			return SerializeLabel(issue.TaskType)
		}()
	}

	return dto
}
