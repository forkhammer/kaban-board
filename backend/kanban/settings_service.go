package kanban

import "main/config"

type SettingsService struct{}

func (s *SettingsService) GetSettings() *KanbanSettings {
	return &KanbanSettings{
		Logo:    config.Settings.Logo,
		Caption: config.Settings.Caption,
	}
}
