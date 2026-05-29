package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"
	"time"

	"github.com/goioc/di"
)

type SprintFactory struct {
	repo repo.SprintRepo
}

func NewSprintFactory() *SprintFactory {
	return &SprintFactory{
		repo: di.GetInstance("SprintRepository").(repo.SprintRepo),
	}
}

func (f *SprintFactory) Build(data map[string]any) *models.Sprint {
	sprint := &models.Sprint{}

	if val, ok := data["title"].(string); ok && val != "" {
		sprint.Title = val
	} else {
		sprint.Title = ""
	}

	if val, ok := data["start_date"].(time.Time); ok {
		sprint.StartDate = val
	} else {
		sprint.StartDate = time.Now().AddDate(0, 0, -7)
	}

	if val, ok := data["end_date"].(time.Time); ok {
		sprint.EndDate = val
	} else {
		sprint.EndDate = time.Now().AddDate(0, 0, 7)
	}

	if val, ok := data["team"].(models.Team); ok {
		sprint.Team = val
	}

	if val, ok := data["hours_per_user"].(uint); ok {
		sprint.HoursPerUser = val
	} else {
		sprint.HoursPerUser = 40
	}

	if val, ok := data["status"].(models.SprintStatus); ok {
		sprint.Status = val
	} else {
		sprint.Status = models.SprintStatusWaiting
	}

	return sprint
}

func (f *SprintFactory) Create(data map[string]any) (*models.Sprint, error) {
	sprint := f.Build(data)
	return f.repo.Create(sprint)
}
