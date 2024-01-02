package kanban

import "main/config"

type ClientSettingsService struct{}

func (s *ClientSettingsService) GetSettings() *KanbanSettings {
	return &KanbanSettings{
		Logo:    config.Settings.Logo,
		Caption: config.Settings.Caption,
	}
}
