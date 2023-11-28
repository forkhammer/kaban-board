package kanban

import (
	"errors"
)

type TeamService struct {
	teamRepository TeamRepository
}

func (s *TeamService) GetAllTeams() ([]Team, error) {
	return s.teamRepository.GetTeams()
}

func (s *TeamService) GetTeamById(id int) (*Team, error) {
	return s.teamRepository.GetTeamById(id)
}

func (s *TeamService) UpdateTeam(id int, data *UpdateTeamRequest) (*Team, error) {
	team, err := s.teamRepository.GetTeamById(id)

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
		return team, nil
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
