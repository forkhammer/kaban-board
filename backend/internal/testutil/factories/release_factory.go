package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type ReleaseFactory struct {
	repo repo.ReleaseRepo
}

func NewReleaseFactory() *ReleaseFactory {
	return &ReleaseFactory{
		repo: di.GetInstance("ReleaseRepository").(repo.ReleaseRepo),
	}
}

func (f *ReleaseFactory) Build(data map[string]any) *models.Release {
	release := &models.Release{}

	if val, ok := data["title"].(string); ok && val != "" {
		release.Title = val
	} else {
		release.Title = gofakeit.AppVersion()
	}

	if val, ok := data["iid"].(string); ok && val != "" {
		release.Iid = models.ReleaseIid(val)
	} else {
		release.Iid = models.ReleaseIid(gofakeit.DigitN(5))
	}

	if val, ok := data["project"].(models.Project); ok {
		release.Project = val
	}

	if val, ok := data["web_path"].(string); ok && val != "" {
		release.WebPath = val
	} else {
		release.WebPath = gofakeit.URL()
	}

	return release
}

func (f *ReleaseFactory) Create(data map[string]any) (*models.Release, error) {
	release := f.Build(data)
	return f.repo.Create(release)
}
