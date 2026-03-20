package account_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type AccountFilterSpec struct {
	repo.QuerySpec
	Filter queries.AccountFilter
}

func (s *AccountFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if len(s.Filter.GitlabIds) > 0 {
		query = query.Where("accounts.gitlab_id IN ?", s.Filter.GitlabIds)
	}

	if len(s.Filter.Ids) > 0 {
		query = query.Where("accounts.id IN ?", s.Filter.Ids)
	}

	return query, nil
}

type AccountQueryImpl struct {
	queries.AccountQuery
}

func (q *AccountQueryImpl) GetSpec(filter queries.AccountFilter) repo.QuerySpec {
	return &AccountFilterSpec{Filter: filter}
}
