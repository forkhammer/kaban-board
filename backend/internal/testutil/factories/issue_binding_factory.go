package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/goioc/di"
)

type IssueBindingFactory struct {
	repo repo.IssueBindingRepo
}

func NewIssueBindingFactory() *IssueBindingFactory {
	return &IssueBindingFactory{
		repo: di.GetInstance("IssueBindingRepository").(repo.IssueBindingRepo),
	}
}

func (f *IssueBindingFactory) Build(data map[string]any) *models.IssueBinding {
	binding := &models.IssueBinding{}

	if val, ok := data["sprint"].(*models.Sprint); ok {
		binding.Sprint = val
	}

	if val, ok := data["issue"].(*models.Issue); ok {
		binding.Issue = val
	}

	if val, ok := data["estimate_dev"].(*uint); ok {
		binding.EstimateDev = val
	}

	if val, ok := data["estimate_qa"].(*uint); ok {
		binding.EstimateQA = val
	}

	if val, ok := data["bind_status"].(models.IssueBindingStatus); ok && val != "" {
		binding.BindStatus = val
	} else {
		binding.BindStatus = models.IssueBindStatusBacklog
	}

	if val, ok := data["priority"].(*models.IssueBindingPriority); ok {
		binding.Priority = val
	} else {
		priority := models.IssueBindingPriorityMedium
		binding.Priority = &priority
	}

	if val, ok := data["assignee"].(*models.User); ok {
		binding.Assignee = val
	}

	if val, ok := data["comment"].(*string); ok {
		binding.Comment = val
	}

	if val, ok := data["release"].(*models.Release); ok {
		binding.Release = val
	}

	if val, ok := data["epic"].(*models.Epic); ok {
		binding.Epic = val
	}

	if val, ok := data["order"].(string); ok {
		binding.Order = val
	}

	if val, ok := data["planned"].(*bool); ok {
		binding.Planned = val
	}

	if val, ok := data["version"].(uint); ok {
		binding.Version = val
	} else {
		binding.Version = 1
	}

	return binding
}

func (f *IssueBindingFactory) Create(data map[string]any) (*models.IssueBinding, error) {
	binding := f.Build(data)
	return f.repo.Create(binding)
}
