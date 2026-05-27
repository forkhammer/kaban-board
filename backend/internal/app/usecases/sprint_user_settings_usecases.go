package usecases

import (
	"errors"
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type CreateSprintUserSettingsRequest struct {
	UserId       uint
	HoursPerUser uint
}

type SprintUserSettingsUseCases struct {
	repo  repo.SprintUserSettingsRepo     `di.inject:"SprintUserSettingsRepository"`
	query queries.SprintUserSettingsQuery `di.inject:"SprintUserSettingsQuery"`
}

func (uc *SprintUserSettingsUseCases) GetSprintUserSettings(filter *queries.SprintUserSettingsFilter) ([]domain.SprintUserSettings, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.query.GetSpec(*filter)
	}
	return uc.repo.List(spec)
}

func (uc *SprintUserSettingsUseCases) CreateOrUpdateSprintUserSettings(sprintId uint, request *CreateSprintUserSettingsRequest) (*domain.SprintUserSettings, error) {
	settingsList, err := uc.repo.List(uc.query.GetSpec(queries.SprintUserSettingsFilter{
		SprintId: &sprintId,
		UserId:   &request.UserId,
	}))
	if err != nil {
		return nil, err
	}

	if len(settingsList) > 0 {
		settings := &settingsList[0]
		settings.HoursPerUser = request.HoursPerUser
		if err := settings.Validate(); err != nil {
			return nil, err
		}
		return uc.repo.Update(settings)
	}

	settings := &domain.SprintUserSettings{
		SprintId:     sprintId,
		UserId:       request.UserId,
		HoursPerUser: request.HoursPerUser,
	}
	if err := settings.Validate(); err != nil {
		return nil, err
	}
	return uc.repo.Create(settings)
}

func (uc *SprintUserSettingsUseCases) DeleteSprintUserSettings(sprintId uint, userId uint) error {
	settingsList, err := uc.repo.List(uc.query.GetSpec(queries.SprintUserSettingsFilter{
		SprintId: &sprintId,
		UserId:   &userId,
	}))
	if err != nil {
		return err
	}

	if len(settingsList) == 0 {
		return errors.New("settings not found")
	}

	return uc.repo.Delete(settingsList[0].Id)
}
