package dto

import (
	domain "main/internal/domain/models"
)

type TeamDto struct {
	Id     int        `json:"id"`
	Title  string     `json:"title"`
	Groups []GroupDto `json:"groups"`
}

func SerializeTeam(team *domain.Team) *TeamDto {
	return &TeamDto{
		Id:     int(team.Id),
		Title:  team.Title,
		Groups: *SerializeGroups(&team.Groups),
	}
}
