package repo

import (
	"errors"
	"strconv"

	"main/internal/app/queries"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	conn         interfaces.ConnectionInterface `di.inject:"db"`
	groupRepo    *GroupRepository               `di.inject:"GroupRepository"`
	accountRepo  *AccountRepository             `di.inject:"AccountRepository"`
	accountQuery queries.AccountQuery           `di.inject:"AccountQuery"`
}

func (r *UserRepository) Get(id domain.UserId) (*domain.User, error) {
	user := &models.User{}
	if err := r.getQuery().Where("users.id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_pkg.NewNotFoundError("User", strconv.Itoa(int(id)), err)
		}
		return nil, err
	}

	accounts, err := r.loadAccountsMap([]uint{user.Id})
	if err != nil {
		return nil, err
	}

	return r.toDomainUser(user, accounts), nil
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

	userIds := make([]uint, len(users))
	for i, user := range users {
		userIds[i] = user.Id
	}

	accounts, err := r.loadAccountsMap(userIds)
	if err != nil {
		return nil, err
	}

	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = *r.toDomainUser(&user, accounts)
	}

	return domainUsers, nil
}

func (r *UserRepository) Count(spec repo.QuerySpec) (int, error) {
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

	return int(count), nil
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

func (r *UserRepository) toDomainUser(user *models.User, accounts map[uint]*domain.Account) *domain.User {
	domainUser := &domain.User{
		Id:        domain.UserId(user.Id),
		IsVisible: user.IsVisible,
		IsActive:  user.IsActive,
		Name:      user.Name,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		Groups: utils.Map(user.Groups, func(group *models.Group) domain.Group {
			return *r.groupRepo.toDomainGroup(group)
		}),
	}

	if account, ok := accounts[user.Id]; ok {
		domainUser.Account = &domain.UserAccount{
			Id:   account.Id,
			Name: account.Name,
			Role: account.Role,
		}
	}

	return domainUser
}

func (r *UserRepository) toUser(user *domain.User) *models.User {
	return &models.User{
		Id:        uint(user.Id),
		IsVisible: user.IsVisible,
		IsActive:  user.IsActive,
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

func (r *UserRepository) loadAccountsMap(userIds []uint) (map[uint]*domain.Account, error) {
	if len(userIds) == 0 {
		return make(map[uint]*domain.Account), nil
	}

	accounts, err := r.accountRepo.List(r.accountQuery.GetSpec(queries.AccountFilter{GitlabIds: userIds}))
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*domain.Account, len(accounts))
	for i := range accounts {
		if accounts[i].GitlabID != nil {
			result[uint(*accounts[i].GitlabID)] = &accounts[i]
		}
	}

	return result, nil
}
