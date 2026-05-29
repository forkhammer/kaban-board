package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type ProjectFactory struct {
	repo repo.ProjectRepo
}

func NewProjectFactory() *ProjectFactory {
	return &ProjectFactory{
		repo: di.GetInstance("ProjectRepository").(repo.ProjectRepo),
	}
}

func (f *ProjectFactory) Build(data map[string]any) *models.Project {
	project := &models.Project{}

	if val, ok := data["name"].(string); ok && val != "" {
		project.Name = val
	} else {
		project.Name = gofakeit.AppName()
	}

	if val, ok := data["is_visible"].(bool); ok {
		project.IsVisible = val
	} else {
		project.IsVisible = true
	}

	if val, ok := data["team"].(*models.Team); ok {
		project.Team = val
	}

	return project
}

func (f *ProjectFactory) Create(data map[string]any) (*models.Project, error) {
	project := f.Build(data)
	return f.repo.Create(project)
}
