package dto

import (
	"fmt"
	domain "main/internal/domain/models"
	"main/pkg/utils"
	"time"
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
}

type KanbanUserDto struct {
	Id        uint       `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	AvatarUrl string     `json:"avatarUrl"`
	Issues    []IssueDto `json:"issues"`
	Teams     []uint     `json:"teams"`
	Groups    []GroupDto `json:"groups"`
}

type BoardDto struct {
	Users      []KanbanUserDto `json:"users"`
	UpdateTime *time.Time      `json:"updateTime"`
}

func SerializeKanbanUser(user *domain.KanbanUser) *KanbanUserDto {
	return &KanbanUserDto{
		Id:        uint(user.User.Id),
		Name:      user.User.Name,
		Username:  user.User.Username,
		AvatarUrl: cleanUserAvatar(user.User.AvatarUrl),
		Issues:    SerializeIssues(user.Issues),
		Teams: utils.Map(user.Teams, func(team domain.Team) uint {
			return uint(team.Id)
		}),
		Groups: SerializeGroups(user.User.Groups),
	}
}

func SerializeKanbanUsers(users []domain.KanbanUser) []KanbanUserDto {
	return utils.Map(users, func(user domain.KanbanUser) KanbanUserDto {
		return *SerializeKanbanUser(&user)
	})
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
	}
}
