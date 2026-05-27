package repo

import (
	"errors"
	"strconv"

	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type SprintUserSettingsRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *SprintUserSettingsRepository) Get(id domain.SprintUserSettingsId) (*domain.SprintUserSettings, error) {
	settings := &models.SprintUserSettings{}
	if err := r.conn.GetEngine().Joins("User").Where("sprint_user_settings.id = ?", id).First(settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_pkg.NewNotFoundError("SprintUserSettings", strconv.Itoa(int(id)), err)
		}
		return nil, err
	}
	return r.toDomainSprintUserSettings(settings), nil
}

func (r *SprintUserSettingsRepository) List(spec repo.QuerySpec) ([]domain.SprintUserSettings, error) {
	settings := make([]models.SprintUserSettings, 0)
	query := r.conn.GetEngine().Model(&models.SprintUserSettings{}).Joins("User")

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&settings).Error; err != nil {
		return nil, err
	}

	domainSettings := make([]domain.SprintUserSettings, len(settings))
	for i, s := range settings {
		domainSettings[i] = *r.toDomainSprintUserSettings(&s)
	}

	return domainSettings, nil
}

func (r *SprintUserSettingsRepository) Count(spec repo.QuerySpec) (int, error) {
	var count int64
	query := r.conn.GetEngine().Model(&models.SprintUserSettings{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return 0, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *SprintUserSettingsRepository) Create(settings *domain.SprintUserSettings) (*domain.SprintUserSettings, error) {
	model := r.toSprintUserSettings(settings)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.SprintUserSettingsId(model.Id))
}

func (r *SprintUserSettingsRepository) Update(settings *domain.SprintUserSettings) (*domain.SprintUserSettings, error) {
	model := r.toSprintUserSettings(settings)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.SprintUserSettingsId(model.Id))
}

func (r *SprintUserSettingsRepository) Delete(id domain.SprintUserSettingsId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.SprintUserSettings{}).Error
}

func (r *SprintUserSettingsRepository) toDomainSprintUserSettings(settings *models.SprintUserSettings) *domain.SprintUserSettings {
	return &domain.SprintUserSettings{
		Id:           domain.SprintUserSettingsId(settings.Id),
		SprintId:     settings.SprintId,
		UserId:       settings.UserId,
		User: domain.User{
			Id:        domain.UserId(settings.User.Id),
			Name:      settings.User.Name,
			AvatarUrl: settings.User.AvatarUrl,
		},
		HoursPerUser: settings.HoursPerUser,
	}
}

func (r *SprintUserSettingsRepository) toSprintUserSettings(settings *domain.SprintUserSettings) *models.SprintUserSettings {
	return &models.SprintUserSettings{
		Id:           uint(settings.Id),
		SprintId:     settings.SprintId,
		UserId:       settings.UserId,
		HoursPerUser: settings.HoursPerUser,
	}
}
