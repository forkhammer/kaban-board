package repo

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"
)

const TASK_TYPE_LABELS_KEY = "task_type_labels"

type SettingsRepository struct {
	repo.SettingsRepo
	kvRepo *KeyValueRepository `di.inject:"KeyValueRepository"`
}

func (r *SettingsRepository) Get() (*models.Settings, error) {
	taskTypeLabels := make([]string, 0)
	err := r.kvRepo.Get(TASK_TYPE_LABELS_KEY, &taskTypeLabels, []string{})
	if err != nil {
		return nil, err
	}
	return &models.Settings{TaskTypeLabels: taskTypeLabels}, nil
}

func (r *SettingsRepository) Save(settings *models.Settings) (*models.Settings, error) {
	err := r.kvRepo.Set(TASK_TYPE_LABELS_KEY, settings.TaskTypeLabels)
	if err != nil {
		return nil, err
	}
	return settings, nil
}
