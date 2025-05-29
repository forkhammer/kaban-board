package column_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type ColumnFilterSpec struct {
	repo.QuerySpec
	Filter queries.ColumnFilter
}

func (s *ColumnFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if len(s.Filter.Ids) > 0 {
		query = query.Where("id IN ?", s.Filter.Ids)
	}

	return query, nil
}

type ColumnQueryImpl struct {
	queries.ColumnQuery
}

func (b *ColumnQueryImpl) GetSpec(filter queries.ColumnFilter) repo.QuerySpec {
	return &ColumnFilterSpec{Filter: filter}
}
