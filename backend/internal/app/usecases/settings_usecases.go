package usecases

import (
	"main/config"
	"main/internal/domain/repo"
)

type ClientSettings struct {
	Logo    string `json:"logo"`
	Caption string `json:"caption"`
}

type SettingsUseCases struct {
	config       *config.Config    `di.inject:"config"`
	settingsRepo repo.SettingsRepo `di.inject:"SettingsRepository"`
}

type KanbanSettings struct {
	TaskTypeLabels []string `json:"taskTypeLabels"`
}

func (s *SettingsUseCases) GetClientSettings() *ClientSettings {
	return &ClientSettings{
		Logo:    s.config.Logo,
		Caption: s.config.Caption,
	}
}

func (s *SettingsUseCases) GetKanbantSettings() (*KanbanSettings, error) {
	settings, err := s.settingsRepo.Get()
	if err != nil {
		return nil, err
	}
	return &KanbanSettings{
		TaskTypeLabels: settings.TaskTypeLabels,
	}, nil
}

func (s *SettingsUseCases) SetTaskTypeLabels(labels []string) (*KanbanSettings, error) {
	settings, err := s.settingsRepo.Get()
	if err != nil {
		return nil, err
	}

	settings.TaskTypeLabels = labels
	settings, err = s.settingsRepo.Save(settings)
	if err != nil {
		return nil, err
	}
	return &KanbanSettings{
		TaskTypeLabels: settings.TaskTypeLabels,
	}, nil
}
