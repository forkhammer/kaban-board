package label_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type LabelFilterSpec struct {
	repo.QuerySpec
	Filter queries.LabelFilter
}

func (s *LabelFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if len(s.Filter.Names) > 0 {
		query = query.Where("name IN ?", s.Filter.Names)
	}

	if s.Filter.Priority != nil {
		query = query.Where("priority = ?", string(*s.Filter.Priority))
	}

	return query, nil
}

type LabelQueryImpl struct {
	queries.LabelQuery
}

func (b *LabelQueryImpl) GetSpec(filter queries.LabelFilter) repo.QuerySpec {
	return &LabelFilterSpec{Filter: filter}
}
