package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
	"time"
)

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
		AvatarUrl: user.User.AvatarUrl,
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
