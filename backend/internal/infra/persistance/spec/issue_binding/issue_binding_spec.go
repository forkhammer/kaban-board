package issuebinding_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type IssueBindingFilterSpec struct {
	repo.QuerySpec
	Filter queries.IssueBindingFilter
}

func (s *IssueBindingFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if len(s.Filter.Issues) > 0 {
		query = query.Where("id in ?", s.Filter.Issues)
	}

	return query, nil
}

type IssueBindingQueryImpl struct {
	queries.IssueBindingQuery
}

func (b *IssueBindingQueryImpl) GetSpec(filter queries.IssueBindingFilter) repo.QuerySpec {
	return &IssueBindingFilterSpec{Filter: filter}
}
