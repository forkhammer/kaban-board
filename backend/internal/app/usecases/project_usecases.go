package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ProjectUseCases struct {
	projectRepo  repo.ProjectRepo     `di.inject:"ProjectRepository"`
	teamRepo     repo.TeamRepo        `di.inject:"TeamRepository"`
	projectQuery queries.ProjectQuery `di.inject:"ProjectQuery"`
}

func (uc *ProjectUseCases) GetProjects(filter *queries.ProjectFilter) ([]domain.Project, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.projectQuery.GetSpec(*filter)
	}

	return uc.projectRepo.List(spec)
}

func (uc *ProjectUseCases) GetProject(id uint) (*domain.Project, error) {
	return uc.projectRepo.Get(domain.ProjectId(id))
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
