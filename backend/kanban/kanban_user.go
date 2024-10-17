package kanban

import "main/repository/models"

type KanbanUser struct {
	Id        uint            `json:"id"`
	Name      string          `json:"name"`
	Username  string          `json:"username"`
	AvatarUrl string          `json:"avatarUrl"`
	Issues    []Issue         `json:"issues"`
	Teams     []uint          `json:"teams"`
	Groups    []*models.Group `json:"groups"`
}
