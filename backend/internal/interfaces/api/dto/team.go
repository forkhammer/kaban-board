package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type TeamDto struct {
	Id     int        `json:"id"`
	Title  string     `json:"title"`
	Groups []GroupDto `json:"groups"`
}

type UpdateTeamRequest struct {
	Title  string `json:"title"`
	Groups []int  `json:"groups"`
}

type CreateTeamRequest struct {
	Title  string `json:"title"`
	Groups []int  `json:"groups"`
}

func SerializeTeam(team *domain.Team) *TeamDto {
	return &TeamDto{
		Id:     int(team.Id),
		Title:  team.Title,
		Groups: *SerializeGroups(&team.Groups),
	}
}

func SerializeTeams(teams *[]domain.Team) *[]TeamDto {
	result := utils.Map(*teams, func(team domain.Team) TeamDto {
		return *SerializeTeam(&team)
	})
	return &result
}
