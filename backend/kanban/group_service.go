package kanban

import (
	"main/repository"
	"main/repository/models"
)

type GroupService struct {
	groupRepository repository.GroupRepositoryInterface `di.inject:"groupRepository"`
}

func (s *GroupService) GetGroups() ([]models.Group, error) {
	var groups []models.Group
	err := s.groupRepository.GetGroups(&groups)
	return groups, err
}

func (s *GroupService) GetGroupById(id int) (*models.Group, error) {
	var group models.Group
	err := s.groupRepository.GetGroupById(&group, id)
	return &group, err
}
