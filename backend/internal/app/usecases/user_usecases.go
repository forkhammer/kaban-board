package usecases

import (
	"errors"
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type UserUseCases struct {
	userRepo    repo.UserRepo      `di.inject:"UserRepository"`
	groupRepo   repo.GroupRepo     `di.inject:"GroupRepository"`
	groupQuery  queries.GroupQuery `di.inject:"GroupQuery"`
	userQuery   queries.UserQuery  `di.inject:"UserQuery"`
	accountRepo repo.AccountRepo   `di.inject:"AccountRepository"`
}

func (uc *UserUseCases) GetUsers(filter *queries.UserFilter) ([]domain.User, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.userQuery.GetSpec(*filter)
	}
	return uc.userRepo.List(spec)
}

func (uc *UserUseCases) GetUser(id uint) (*domain.User, error) {
	return uc.userRepo.Get(domain.UserId(id))
}

func (uc *UserUseCases) SetVisibility(id uint, visible bool) (*domain.User, error) {
	user, err := uc.userRepo.Get(domain.UserId(id))
	if err != nil {
		return nil, err
	}

	user.IsVisible = visible

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return uc.userRepo.Update(user)
}

func (uc *UserUseCases) SetGroups(id uint, groupIds []uint) (*domain.User, error) {
	user, err := uc.userRepo.Get(domain.UserId(id))
	if err != nil {
		return nil, err
	}

	groups, err := uc.groupRepo.List(uc.groupQuery.GetSpec(queries.GroupFilter{Ids: utils.Map(groupIds, func(g uint) domain.GroupId {
		return domain.GroupId(g)
	})}))
	if err != nil {
		return nil, err
	}
	user.Groups = groups

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return uc.userRepo.Update(user)
}

func (uc *UserUseCases) SetAccountRole(userId uint, role domain.AccountRole) (*domain.User, error) {
	account, err := uc.accountRepo.GetByGitlabID(userId)
	if err != nil {
		return nil, errors.New("С пользователем не связан аккаунт")
	}

	account.Role = role

	_, err = uc.accountRepo.Update(account)
	if err != nil {
		return nil, err
	}

	return uc.userRepo.Get(domain.UserId(userId))
}
