package group_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type GroupFilterSpec struct {
	repo.QuerySpec
	Filter queries.GroupFilter
}

func (s *GroupFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if len(s.Filter.Ids) > 0 {
		query = query.Where("id IN ?", s.Filter.Ids)
	}

	return query, nil
}

type GroupQueryImpl struct {
	queries.GroupQuery
}

func (b *GroupQueryImpl) GetSpec(filter queries.GroupFilter) repo.QuerySpec {
	return &GroupFilterSpec{Filter: filter}
}
