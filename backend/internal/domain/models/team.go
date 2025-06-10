package models

import "fmt"

type TeamId uint

type Team struct {
	Id     TeamId
	Title  string
	Groups []Group
}

func NewTeam(id TeamId, title string, groups []Group) (*Team, error) {
	team := &Team{
		Id:     id,
		Title:  title,
		Groups: groups,
	}
	err := team.Validate()
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (t *Team) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("Team title is required")
	}

	return nil
}
