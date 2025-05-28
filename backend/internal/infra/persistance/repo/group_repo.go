package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type GroupRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *GroupRepository) Get(id domain.GroupId) (*domain.Group, error) {
	group := &models.Group{}
	if err := r.conn.GetEngine().Where("id = ?", id).First(group).Error; err != nil {
		return nil, err
	}
	return r.toDomainGroup(group), nil
}

func (r *GroupRepository) List(spec repo.QuerySpec) (*[]domain.Group, error) {
	groups := make([]models.Group, 0)
	query := r.conn.GetEngine().Model(&models.Group{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}

	domainGroups := make([]domain.Group, len(groups))
	for i, group := range groups {
		domainGroups[i] = *r.toDomainGroup(&group)
	}

	return &domainGroups, nil
}

func (r *GroupRepository) Create(group *domain.Group) error {
	model := r.toGroup(group)
	return r.conn.GetEngine().Create(model).Error
}

func (r *GroupRepository) Update(group *domain.Group) error {
	model := r.toGroup(group)
	return r.conn.GetEngine().Save(model).Error
}

func (r *GroupRepository) Delete(id domain.GroupId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Group{}).Error
}

func (r *GroupRepository) toDomainGroup(group *models.Group) *domain.Group {
	return &domain.Group{
		Id:   domain.GroupId(group.Id),
		Name: group.Name,
	}
}

func (r *GroupRepository) toGroup(group *domain.Group) *models.Group {
	return &models.Group{
		Id:   uint(group.Id),
		Name: group.Name,
	}
}
