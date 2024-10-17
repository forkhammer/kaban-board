package kanban

import (
	"main/tools"
)

type GroupService struct {
	groupRepository tools.GroupRepositoryInterface `di.inject:"groupRepository"`
}

func (s *GroupService) GetGroups() ([]Group, error) {
	var groups []Group
	err := s.groupRepository.GetGroups(&groups)
	return groups, err
}

func (s *GroupService) GetGroupById(id int) (*Group, error) {
	var group Group
	err := s.groupRepository.GetGroupById(&group, id)
	return &group, err
}
