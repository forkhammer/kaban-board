package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type EpicFactory struct {
	repo repo.EpicRepo
}

func NewEpicFactory() *EpicFactory {
	return &EpicFactory{
		repo: di.GetInstance("EpicRepository").(repo.EpicRepo),
	}
}

func (f *EpicFactory) Build(data map[string]any) *models.Epic {
	epic := &models.Epic{}

	if val, ok := data["title"].(string); ok && val != "" {
		epic.Title = val
	} else {
		epic.Title = gofakeit.Sentence(3)
	}

	if val, ok := data["external_id"].(string); ok && val != "" {
		epic.ExternalId = models.IssueExternalId(val)
	} else {
		epic.ExternalId = models.IssueExternalId(gofakeit.UUID())
	}

	if val, ok := data["iid"].(string); ok && val != "" {
		epic.Iid = models.IssueIid(val)
	} else {
		epic.Iid = models.IssueIid(gofakeit.DigitN(5))
	}

	if val, ok := data["project"].(models.Project); ok {
		epic.Project = val
	}

	return epic
}

func (f *EpicFactory) Create(data map[string]any) (*models.Epic, error) {
	epic := f.Build(data)
	return f.repo.Create(epic)
}
