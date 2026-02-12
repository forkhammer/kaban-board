package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	conn      interfaces.ConnectionInterface `di.inject:"db"`
	groupRepo *GroupRepository               `di.inject:"GroupRepository"`
}

func (r *UserRepository) Get(id domain.UserId) (*domain.User, error) {
	user := &models.User{}
	if err := r.getQuery().Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return r.toDomainUser(user), nil
}

func (r *UserRepository) List(spec repo.QuerySpec) ([]domain.User, error) {
	users := make([]models.User, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = *r.toDomainUser(&user)
	}

	return domainUsers, nil
}

func (r *UserRepository) Count(spec repo.QuerySpec) (int64, error) {
	var count int64
	query := r.getQuery()

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

	return count, nil
}

func (r *UserRepository) Create(user *domain.User) (*domain.User, error) {
	model := r.toUser(user)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}

	err = r.conn.GetEngine().Model(model).Association("Groups").Replace(model.Groups)
	if err != nil {
		return nil, err
	}

	return r.Get(domain.UserId(model.Id))
}

func (r *UserRepository) Update(user *domain.User) (*domain.User, error) {
	model := r.toUser(user)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}

	err = r.conn.GetEngine().Model(model).Association("Groups").Replace(model.Groups)
	if err != nil {
		return nil, err
	}

	return r.Get(domain.UserId(model.Id))
}

func (r *UserRepository) Delete(id domain.UserId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *UserRepository) toDomainUser(user *models.User) *domain.User {
	return &domain.User{
		Id:        domain.UserId(user.Id),
		IsVisible: user.IsVisible,
		Name:      user.Name,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		Groups: utils.Map(user.Groups, func(group *models.Group) domain.Group {
			return *r.groupRepo.toDomainGroup(group)
		}),
	}
}

func (r *UserRepository) toUser(user *domain.User) *models.User {
	return &models.User{
		Id:        uint(user.Id),
		IsVisible: user.IsVisible,
		Name:      user.Name,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		Groups: utils.Map(user.Groups, func(group domain.Group) *models.Group {
			return r.groupRepo.toGroup(&group)
		}),
	}
}

func (r *UserRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.User{}).Preload("Groups")
}
