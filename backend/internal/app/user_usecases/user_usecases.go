package user_usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type UserUseCases struct {
	userRepo   repo.UserRepo      `di.inject:"UserRepository"`
	groupRepo  repo.GroupRepo     `di.inject:"GroupRepository"`
	groupQuery queries.GroupQuery `di.inject:"GroupQuery"`
}

func (uc *UserUseCases) GetUsers() (*[]domain.User, error) {
	return uc.userRepo.List(nil)
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
	user.Groups = *groups

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return uc.userRepo.Update(user)
}
