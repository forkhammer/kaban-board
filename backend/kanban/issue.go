package kanban

import (
	"main/gitlab"
	"main/repository/models"
)

type Issue struct {
	Id          string                  `json:"id"`
	Iid         string                  `json:"iid"`
	Title       string                  `json:"title"`
	IssueType   string                  `json:"type"`
	Assignees   []gitlab.GitlabAssignee `json:"assignees"`
	WebUrl      string                  `json:"webUrl"`
	Labels      []models.Label          `json:"labels"`
	ProjectId   int                     `json:"projectId"`
	ProjectName *string                 `json:"projectName"`
	Milestone   gitlab.GitlabMilestone  `json:"milestone"`
	TaskType    *models.Label           `json:"taskType"`
}
