package kanban

import (
	"main/config"
	"main/repository/models"
)

type ClientSettingsService struct{}

func (s *ClientSettingsService) GetSettings() *models.ClientSettings {
	return &models.ClientSettings{
		Logo:    config.Settings.Logo,
		Caption: config.Settings.Caption,
	}
}
