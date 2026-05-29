package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type TeamFactory struct {
	repo repo.TeamRepo
}

func NewTeamFactory() *TeamFactory {
	return &TeamFactory{
		repo: di.GetInstance("TeamRepository").(repo.TeamRepo),
	}
}

func (f *TeamFactory) Build(data map[string]any) *models.Team {
	team := &models.Team{}

	if val, ok := data["title"].(string); ok && val != "" {
		team.Title = val
	} else {
		team.Title = gofakeit.Company()
	}

	if val, ok := data["groups"].([]models.Group); ok {
		team.Groups = val
	}

	return team
}

func (f *TeamFactory) Create(data map[string]any) (*models.Team, error) {
	team := f.Build(data)
	return f.repo.Create(team)
}
