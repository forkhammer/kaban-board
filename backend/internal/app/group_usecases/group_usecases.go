package group_usecases

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type GroupUseCases struct {
	groupRepo repo.GroupRepo `di.inject:"GroupRepository"`
}

func (uc *GroupUseCases) GetGroups() (*[]domain.Group, error) {
	return uc.groupRepo.List(nil)
}

func (uc *GroupUseCases) GetGroup(id uint) (*domain.Group, error) {
	return uc.groupRepo.Get(domain.GroupId(id))
}
