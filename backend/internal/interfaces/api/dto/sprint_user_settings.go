package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type SprintUserSettingsUserDto struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type SprintUserSettingsDto struct {
	Id           uint                       `json:"id"`
	SprintId     uint                       `json:"sprint_id"`
	UserId       uint                       `json:"user_id"`
	User         SprintUserSettingsUserDto  `json:"user"`
	HoursPerUser uint                       `json:"hours_per_user"`
}

type CreateSprintUserSettingsRequest struct {
	UserId       uint `json:"user_id" binding:"required"`
	HoursPerUser uint `json:"hours_per_user" binding:"required"`
}

func SerializeSprintUserSettings(settings *domain.SprintUserSettings) *SprintUserSettingsDto {
	return &SprintUserSettingsDto{
		Id:       uint(settings.Id),
		SprintId: settings.SprintId,
		UserId:   settings.UserId,
		User: SprintUserSettingsUserDto{
			Id:        uint(settings.User.Id),
			Name:      settings.User.Name,
			AvatarUrl: settings.User.AvatarUrl,
		},
		HoursPerUser: settings.HoursPerUser,
	}
}

func SerializeSprintUserSettingsList(settings []domain.SprintUserSettings) []SprintUserSettingsDto {
	return utils.Map(settings, func(s domain.SprintUserSettings) SprintUserSettingsDto {
		return *SerializeSprintUserSettings(&s)
	})
}
