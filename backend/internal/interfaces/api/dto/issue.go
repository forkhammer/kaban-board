package dto

import (
	"fmt"
	"main/internal/app/usecases"
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
	Page     int     `form:"page,default=1"`
	Limit    int     `form:"limit,default=0"`
}

type BindIssueRequest struct {
	IssueIds   []uint `json:"issue_ids"`
	SprintId   uint   `json:"sprint_id"`
	AssigneeId uint   `json:"assignee_id"`
}

type BindIssueResponse struct {
	Results []IssueBindingDto `json:"results"`
	Errors  []string          `json:"errors"`
}

type SaveIssueBindingRequest struct {
	Title       *string `json:"title"`
	EstimateDev *uint   `json:"estimateDev"`
	EstimateQA  *uint   `json:"estimateQA"`
	BindStatus  *string `json:"bindStatus"`
	Assignee    *uint   `json:"assignee"`
	Comment     *string `json:"comment"`
	Priority    *string `json:"priority"`
	Release     *uint   `json:"release"`
	Epic        *uint  `json:"epic"`
	Planned     *bool  `json:"planned"`
	Version     uint   `json:"version"`
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
	Release     *ReleaseDto `json:"release"`
	Epic        *EpicDto    `json:"epic"`
	TaskType    *LabelDto   `json:"taskType"`
	EstimateDev *uint       `json:"estimateDev"`
	EstimateQA  *uint       `json:"estimateQA"`
	BindingId   *uint       `json:"bindingId"`
	BindStatus  *string     `json:"bindStatus"`
	Priority    *string     `json:"priority"`
	Comment     *string     `json:"comment"`
}

type IssuePageDto struct {
	Count      int        `json:"count"`
	Pages      int        `json:"pages"`
	Page       int        `json:"page"`
	StartIndex int        `json:"start_index"`
	EndIndex   int        `json:"end_index"`
	Results    []IssueDto `json:"results"`
}

func SerializeIssuePage(issuePage *usecases.IssuePage) *IssuePageDto {
	totalPages := int(issuePage.Count / issuePage.Limit)
	if issuePage.Count%issuePage.Limit > 0 {
		totalPages++
	}

	startIndex := int((issuePage.Page - 1) * issuePage.Limit)
	endIndex := startIndex + len(issuePage.Results)

	return &IssuePageDto{
		Count:      int(issuePage.Count),
		Pages:      totalPages,
		Page:       int(issuePage.Page),
		StartIndex: startIndex,
		EndIndex:   endIndex,
		Results:    SerializeIssues(issuePage.Results),
	}
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
		Labels:      SerializeLabels(issue.Labels),
		ProjectId:   int(issue.Project.Id),
		ProjectName: &issue.Project.Name,
		Release: func() *ReleaseDto {
			release := issue.Release
			if release == nil {
				return nil
			}
			return SerializeRelease(release)
		}(),
		Epic: nil,
		TaskType: func() *LabelDto {
			if issue.TaskType == nil {
				return nil
			}
			return SerializeLabel(issue.TaskType)
		}(),
		EstimateDev: nil,
		EstimateQA:  nil,
		BindingId:   nil,
		BindStatus:  nil,
		Priority:    nil,
		Comment:     nil,
	}
}
