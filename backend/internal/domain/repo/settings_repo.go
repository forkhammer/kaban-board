package repo

import "main/internal/domain/models"

type SettingsRepo interface {
	Get() (*models.Settings, error)
	Save(settings *models.Settings) (*models.Settings, error)
}
