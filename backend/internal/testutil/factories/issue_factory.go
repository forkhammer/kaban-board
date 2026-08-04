package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type IssueFactory struct {
	repo repo.IssueRepo
}

func NewIssueFactory() *IssueFactory {
	return &IssueFactory{
		repo: di.GetInstance("IssueRepository").(repo.IssueRepo),
	}
}

func (f *IssueFactory) Build(data map[string]any) *models.Issue {
	issue := &models.Issue{}

	if val, ok := data["title"].(string); ok && val != "" {
		issue.Title = val
	} else {
		issue.Title = gofakeit.Sentence(3)
	}

	if val, ok := data["issue_type"].(models.IssueType); ok && val != "" {
		issue.IssueType = val
	} else {
		issue.IssueType = models.IssueTypeTask
	}

	if val, ok := data["external_id"].(models.IssueExternalId); ok && val != "" {
		issue.ExternalId = val
	} else {
		issue.ExternalId = models.IssueExternalId(gofakeit.UUID())
	}

	if val, ok := data["iid"].(models.IssueIid); ok && val != "" {
		issue.Iid = val
	} else {
		issue.Iid = models.IssueIid(gofakeit.UUID())
	}

	if val, ok := data["web_url"].(string); ok && val != "" {
		issue.WebUrl = val
	} else {
		issue.WebUrl = gofakeit.URL()
	}

	if val, ok := data["project"].(models.Project); ok {
		issue.Project = val
	}

	if val, ok := data["release"].(*models.Release); ok {
		issue.Release = val
	}

	if val, ok := data["task_type"].(*models.Label); ok {
		issue.TaskType = val
	}

	if val, ok := data["epic"].(*models.Epic); ok {
		issue.Epic = val
	}

	if val, ok := data["estimate_dev"].(*uint); ok {
		issue.EstimateDev = val
	}

	if val, ok := data["estimate_qa"].(*uint); ok {
		issue.EstimateQA = val
	}

	return issue
}

func (f *IssueFactory) Create(data map[string]any) (*models.Issue, error) {
	issue := f.Build(data)
	return f.repo.Create(issue)
}
