package kanban

import "main/config"

type ClientSettingsService struct{}

func (s *ClientSettingsService) GetSettings() *ClientSettings {
	return &ClientSettings{
		Logo:    config.Settings.Logo,
		Caption: config.Settings.Caption,
	}
}
