package queries

import "main/internal/domain/repo"

type IssueBindingFilter struct {
	Issues []uint
}

type IssueBindingQuery interface {
	GetSpec(filter IssueBindingFilter) repo.QuerySpec
}
