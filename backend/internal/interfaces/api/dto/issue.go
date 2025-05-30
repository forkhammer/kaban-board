package dto

import (
	"fmt"
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type IssueDto struct {
	Id          string      `json:"id"`
	Iid         string      `json:"iid"`
	Title       string      `json:"title"`
	IssueType   string      `json:"type"`
	Assignees   []UserDto   `json:"assignees"`
	WebUrl      string      `json:"webUrl"`
	Labels      []LabelDto  `json:"labels"`
	ProjectId   int         `json:"projectId"`
	ProjectName *string     `json:"projectName"`
	Milestone   *ReleaseDto `json:"milestone"`
	TaskType    *LabelDto   `json:"taskType"`
	EstimateDev *uint       `json:"estimateDev"`
	EstimateQA  *uint       `json:"estimateQA"`
}

func SerializeIssues(issues []domain.Issue) []IssueDto {
	return utils.Map(issues, func(issue domain.Issue) IssueDto {
		return *SerializeIssue(&issue)
	})
}

func SerializeIssue(issue *domain.Issue) *IssueDto {
	return &IssueDto{
		Id:          fmt.Sprintf("%d", issue.Id),
		Iid:         string(issue.Iid),
		Title:       issue.Title,
		IssueType:   string(issue.IssueType),
		Assignees:   SerializeUsers(issue.Assignees),
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
		EstimateDev: issue.EstimateDev,
		EstimateQA:  issue.EstimateQA,
	}
}
