package usecases

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ProjectUseCases struct {
	projectRepo repo.ProjectRepo `di.inject:"ProjectRepository"`
	teamRepo    repo.TeamRepo    `di.inject:"TeamRepository"`
}

func (uc *ProjectUseCases) GetProjects() (*[]domain.Project, error) {
	return uc.projectRepo.List(nil)
}

func (uc *ProjectUseCases) SetTeam(id uint, teamId *uint) (*domain.Project, error) {
	project, err := uc.projectRepo.Get(domain.ProjectId(id))
	if err != nil {
		return nil, err
	}
	var team *domain.Team
	if teamId != nil {
		var err error
		team, err = uc.teamRepo.Get(domain.TeamId(*teamId))
		if err != nil {
			return nil, err
		}
	}
	project.Team = team

	if err := project.Validate(); err != nil {
		return nil, err
	}

	return uc.projectRepo.Update(project)
}
