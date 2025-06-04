package usecases

import (
	"fmt"
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"time"
)

type CreateSprintRequest struct {
	Title        string
	StartDate    time.Time
	EndDate      time.Time
	HoursPerUser uint
	TeamId       uint
}

type UpdateSprintRequest struct {
	Id           uint
	Title        string
	StartDate    time.Time
	EndDate      time.Time
	HoursPerUser uint
	TeamId       uint
}

type SprintUseCases struct {
	sprintRepo  repo.SprintRepo     `di.inject:"SprintRepository"`
	teamRepo    repo.TeamRepo       `di.inject:"TeamRepository"`
	sprintQuery queries.SprintQuery `di.inject:"SprintQuery"`
}

func (uc *SprintUseCases) GetSprints(filter *queries.SprintFilter) (*[]domain.Sprint, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.sprintQuery.GetSpec(*filter)
	}
	return uc.sprintRepo.List(spec)
}

func (uc *SprintUseCases) GetSprint(id uint) (*domain.Sprint, error) {
	return uc.sprintRepo.Get(domain.SprintId(id))
}

func (uc *SprintUseCases) CreateSprint(request *CreateSprintRequest) (*domain.Sprint, error) {
	team, err := uc.teamRepo.Get(domain.TeamId(request.TeamId))
	if err != nil {
		return nil, err
	}

	sprint := &domain.Sprint{
		Title:        request.Title,
		StartDate:    request.StartDate,
		EndDate:      request.EndDate,
		HoursPerUser: request.HoursPerUser,
		Team:         *team,
	}
	if err := sprint.Validate(); err != nil {
		return nil, err
	}
	return uc.sprintRepo.Create(sprint)
}

func (uc *SprintUseCases) UpdateSprint(request *UpdateSprintRequest) (*domain.Sprint, error) {
	sprint, err := uc.sprintRepo.Get(domain.SprintId(request.Id))
	if err != nil {
		return nil, err
	}

	team, err := uc.teamRepo.Get(domain.TeamId(request.TeamId))
	if err != nil {
		return nil, err
	}

	sprint.Title = request.Title
	sprint.StartDate = request.StartDate
	sprint.EndDate = request.EndDate
	sprint.HoursPerUser = request.HoursPerUser
	sprint.Team = *team
	if err := sprint.Validate(); err != nil {
		return nil, err
	}
	return uc.sprintRepo.Update(sprint)
}

func (uc *SprintUseCases) DeleteSprint(id uint) error {
	sprint, err := uc.sprintRepo.Get(domain.SprintId(id))
	if err != nil {
		return err
	}
	if sprint.GetStatus() != domain.SprintStatusWaiting {
		return fmt.Errorf("Нельзя удалить этот спринт")
	}

	return uc.sprintRepo.Delete(domain.SprintId(id))
}

func (uc *SprintUseCases) CompleteSprint(id uint) (*domain.Sprint, error) {
	sprint, err := uc.sprintRepo.Get(domain.SprintId(id))
	if err != nil {
		return nil, err
	}
	err = sprint.SetCompleted()
	if err != nil {
		return nil, err
	}
	return uc.sprintRepo.Update(sprint)
}
