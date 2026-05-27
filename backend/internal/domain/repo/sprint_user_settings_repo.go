package repo

import "main/internal/domain/models"

type SprintUserSettingsRepo interface {
	RWRepo[models.SprintUserSettings, models.SprintUserSettingsId]
}
