package app_services

import (
	"time"

	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type IssueBindingHistoryService struct {
	historyRepo repo.IssueBindingHistoryRepo `di.inject:"IssueBindingHistoryRepository"`
}

func (s *IssueBindingHistoryService) AddHistory(binding *domain.IssueBinding) (*domain.IssueBindingHistory, error) {
	history := &domain.IssueBindingHistory{
		IssueBindingId: binding.Id,
		IssueId:        binding.Issue.Id,
		SprintId:       binding.Sprint.Id,
		EstimateDev:    binding.EstimateDev,
		EstimateQA:     binding.EstimateQA,
		BindStatus:     binding.BindStatus,
		Priority:       binding.Priority,
		Assignee:       binding.Assignee,
		CreatedAt:      time.Now(),
	}

	return s.historyRepo.Create(history)
}

func (s *IssueBindingHistoryService) AddDeleteHistory(binding *domain.IssueBinding, deletedAt time.Time) (*domain.IssueBindingHistory, error) {
	history := &domain.IssueBindingHistory{
		IssueBindingId: binding.Id,
		IssueId:        binding.Issue.Id,
		SprintId:       binding.Sprint.Id,
		EstimateDev:    binding.EstimateDev,
		EstimateQA:     binding.EstimateQA,
		BindStatus:     binding.BindStatus,
		Priority:       binding.Priority,
		Assignee:       binding.Assignee,
		RemovedAt:      &deletedAt,
		CreatedAt:      time.Now(),
	}

	return s.historyRepo.Create(history)
}
