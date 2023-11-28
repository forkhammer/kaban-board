package kanban

import (
	"main/gitlab"
)

type KanbanUser struct {
	Id        uint                 `json:"id"`
	Name      string               `json:"name"`
	Username  string               `json:"username"`
	AvatarUrl string               `json:"avatarUrl"`
	Issues    []gitlab.GitlabIssue `json:"issues"`
	Teams     []uint               `json:"teams"`
}
