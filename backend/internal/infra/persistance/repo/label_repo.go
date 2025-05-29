package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type LabelRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *LabelRepository) Get(id domain.LabelId) (*domain.Label, error) {
	label := &models.Label{}
	if err := r.conn.GetEngine().Where("id = ?", id).First(label).Error; err != nil {
		return nil, err
	}
	return r.toDomainLabel(label), nil
}

func (r *LabelRepository) List(spec repo.QuerySpec) (*[]domain.Label, error) {
	labels := make([]models.Label, 0)
	query := r.conn.GetEngine().Model(&models.Label{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&labels).Error; err != nil {
		return nil, err
	}

	domainLabels := make([]domain.Label, len(labels))
	for i, label := range labels {
		domainLabels[i] = *r.toDomainLabel(&label)
	}

	return &domainLabels, nil
}

func (r *LabelRepository) Create(label *domain.Label) (*domain.Label, error) {
	model := r.toLabel(label)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.LabelId(model.Id))
}

func (r *LabelRepository) Update(label *domain.Label) (*domain.Label, error) {
	model := r.toLabel(label)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.LabelId(model.Id))
}

func (r *LabelRepository) Delete(id domain.LabelId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Label{}).Error
}

func (r *LabelRepository) toDomainLabel(label *models.Label) *domain.Label {
	return &domain.Label{
		Id:        domain.LabelId(label.Id),
		Name:      label.Name,
		Color:     domain.Color(label.Color),
		TextColor: domain.Color(label.TextColor),
		AltName:   label.AltName,
	}
}

func (r *LabelRepository) toLabel(label *domain.Label) *models.Label {
	return &models.Label{
		Id:        string(label.Id),
		Name:      label.Name,
		Color:     string(label.Color),
		TextColor: string(label.TextColor),
		AltName:   label.AltName,
	}
}
