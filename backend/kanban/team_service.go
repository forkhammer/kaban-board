package kanban

import (
	"errors"
	"main/repository"
	"main/repository/models"
)

type TeamService struct {
	teamRepository  repository.TeamRepositoryInterface  `di.inject:"teamRepository"`
	groupRepository repository.GroupRepositoryInterface `di.inject:"groupRepository"`
}

func (s *TeamService) GetAllTeams() ([]models.Team, error) {
	var teams []models.Team
	err := s.teamRepository.GetTeams(&teams)
	return teams, err
}

func (s *TeamService) GetTeamById(id int) (*models.Team, error) {
	var team models.Team
	err := s.teamRepository.GetTeamById(&team, id)
	return &team, err
}

func (s *TeamService) UpdateTeam(id int, data *UpdateTeamRequest) (*models.Team, error) {
	var team models.Team
	err := s.teamRepository.GetTeamById(&team, id)

	if err != nil {
		return nil, err
	}

	if data.Title == "" {
		return nil, errors.New("Название группы не может быть пустым")
	}

	team.Title = data.Title
	var groups []*models.Group
	if err := s.groupRepository.GetGroupsByIds(&groups, data.Groups); err != nil {
		return nil, err
	}
	team.Groups = groups

	if err = s.teamRepository.SaveTeam(&team); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func (s *TeamService) CreateTeam(data *CreateTeamRequest) (*models.Team, error) {
	if data.Title == "" {
		return nil, errors.New("Название колонки не может быть пустым")
	}

	team := models.Team{
		Title: data.Title,
	}

	if err := s.teamRepository.CreateTeam(&team); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func (s *TeamService) DeleteTeamById(id int) error {
	var team models.Team
	if err := s.teamRepository.GetTeamById(&team, id); err != nil {
		return err
	}
	return s.teamRepository.DeleteTeam(&team)
}
