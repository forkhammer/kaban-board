package kanban

import (
	"errors"
	"main/tools"
)

type TeamService struct {
	teamRepository tools.TeamRepositoryInterface `di.inject:"teamRepository"`
}

func (s *TeamService) GetAllTeams() ([]Team, error) {
	var teams []Team
	err := s.teamRepository.GetTeams(&teams)
	return teams, err
}

func (s *TeamService) GetTeamById(id int) (*Team, error) {
	var team Team
	err := s.teamRepository.GetTeamById(&team, id)
	return &team, err
}

func (s *TeamService) UpdateTeam(id int, data *UpdateTeamRequest) (*Team, error) {
	var team Team
	err := s.teamRepository.GetTeamById(&team, id)

	if err != nil {
		return nil, err
	}

	if data.Title == "" {
		return nil, errors.New("Название группы не может быть пустым")
	}

	team.Title = data.Title

	if err = s.teamRepository.SaveTeam(team); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func (s *TeamService) CreateTeam(data *CreateTeamRequest) (*Team, error) {
	if data.Title == "" {
		return nil, errors.New("Название колонки не может быть пустым")
	}

	team := Team{
		Title: data.Title,
	}

	if err := s.teamRepository.CreateTeam(&team); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func (s *TeamService) DeleteTeamById(id int) error {
	return s.teamRepository.DeleteTeam(id)
}
