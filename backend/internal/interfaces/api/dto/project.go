package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type SetTeamRequest struct {
	TeamId *uint `json:"team_id"`
}

type ProjectDto struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	IsVisible bool   `json:"is_visible"`
	TeamId    *uint  `json:"team_id"`
}

func SerializeProject(project *domain.Project) *ProjectDto {
	var teamId *uint
	if project.Team != (*domain.Team)(nil) {
		val := uint(project.Team.Id)
		teamId = &val
	}

	return &ProjectDto{
		Id:        uint(project.Id),
		Name:      project.Name,
		IsVisible: project.IsVisible,
		TeamId:    teamId,
	}
}

func SerializeProjects(projects *[]domain.Project) *[]ProjectDto {
	result := utils.Map(*projects, func(project domain.Project) ProjectDto {
		return *SerializeProject(&project)
	})
	return &result
}
