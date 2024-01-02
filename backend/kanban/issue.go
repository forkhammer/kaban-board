package kanban

import (
	"main/gitlab"
)

type Issue struct {
	Id        string `json:"id"`
	Iid       string `json:"iid"`
	Title     string `json:"title"`
	IssueType string `json:"type"`
	Assignees struct {
		Nodes []gitlab.GitlabAssignee `json:"nodes"`
	} `json:"assignees"`
	WebUrl string `json:"webUrl"`
	Labels struct {
		Nodes []gitlab.GitlabLabel `json:"nodes"`
	} `json:"labels"`
	ProjectId   int                    `json:"projectId"`
	ProjectName *string                `json:"projectName"`
	Milestone   gitlab.GitlabMilestone `json:"milestone"`
	TaskType    *gitlab.GitlabLabel    `json:"taskType"`
}
