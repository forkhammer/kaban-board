package dto

import (
	"fmt"
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type IssuesRequest struct {
	Assignee *uint   `form:"assignee"`
	Team     *uint   `form:"team"`
	Group    *uint   `form:"group"`
	Sprint   *uint   `form:"sprint"`
	Project  *uint   `form:"project"`
	Search   *string `form:"search"`
}

type BindIssueRequest struct {
	SprintId uint `json:"sprint_id"`
}

type UnbindIssueRequest struct {
	BindingId uint `json:"binding_id"`
}

type SaveIssueRequest struct {
	BindingId   *uint   `json:"bindingId"`
	EstimateDev *uint   `json:"estimateDev"`
	EstimateQA  *uint   `json:"estimateQA"`
	BindStatus  *string `json:"bindStatus"`
	Assignee    *uint   `json:"assignee"`
	Comment     *string `json:"comment"`
}

type IssueDto struct {
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
	Milestone   *ReleaseDto `json:"milestone"`
	TaskType    *LabelDto   `json:"taskType"`
	EstimateDev *uint       `json:"estimateDev"`
	EstimateQA  *uint       `json:"estimateQA"`
	BindingId   *uint       `json:"bindingId"`
	BindStatus  *string     `json:"bindStatus"`
	Comment     *string     `json:"comment"`
}

func SerializeIssues(issues []domain.Issue) []IssueDto {
	return utils.Map(issues, func(issue domain.Issue) IssueDto {
		return *SerializeIssue(&issue)
	})
}

func SerializeIssue(issue *domain.Issue) *IssueDto {
	return &IssueDto{
		Id:        fmt.Sprintf("%d", issue.Id),
		Iid:       string(issue.Iid),
		Title:     issue.Title,
		IssueType: string(issue.IssueType),
		Assignees: SerializeUsers(issue.Assignees),
		Assignee: func() *UserDto {
			assignee := issue.GetAssignee()
			if assignee == nil {
				return nil
			}
			return SerializeUser(assignee)
		}(),
		WebUrl:      issue.WebUrl,
		Labels:      *SerializeLabels(&issue.Labels),
		ProjectId:   int(issue.Project.Id),
		ProjectName: &issue.Project.Name,
		Milestone: func() *ReleaseDto {
			if issue.Release == nil {
				return nil
			}
			return SerializeRelease(issue.Release)
		}(),
		TaskType: func() *LabelDto {
			if issue.TaskType == nil {
				return nil
			}
			return SerializeLabel(issue.TaskType)
		}(),
		EstimateDev: issue.GetEstimateDev(),
		EstimateQA:  issue.GetEstimateQA(),
		BindingId:   (*uint)(issue.GetContextBindingId()),
		BindStatus:  (*string)(issue.GetBindStatus()),
		Comment:     issue.GetComment(),
	}
}
