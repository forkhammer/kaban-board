package queries

import (
	"main/internal/domain/repo"
	domain "main/internal/domain/models"
)

type LabelFilter struct {
	Names    []string
	Priority *domain.IssueBindingPriority
}

type LabelQuery interface {
	GetSpec(filter LabelFilter) repo.QuerySpec
}
