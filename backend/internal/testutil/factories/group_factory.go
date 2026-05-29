package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type GroupFactory struct {
	repo repo.GroupRepo
}

func NewGroupFactory() *GroupFactory {
	return &GroupFactory{
		repo: di.GetInstance("GroupRepository").(repo.GroupRepo),
	}
}

func (f *GroupFactory) Build(data map[string]any) *models.Group {
	group := &models.Group{}

	if val, ok := data["name"].(string); ok && val != "" {
		group.Name = val
	} else {
		group.Name = gofakeit.Company()
	}

	return group
}

func (f *GroupFactory) Create(data map[string]any) (*models.Group, error) {
	group := f.Build(data)
	return f.repo.Create(group)
}
