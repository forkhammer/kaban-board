package spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type PaginationSpec struct {
	repo.BaseQuerySpec
	page  int
	limit int
}

func (s *PaginationSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.limit > 0 {
		query = query.Limit(s.limit)

		if s.page > 0 {
			offset := (s.page - 1) * s.limit
			query = query.Offset(offset)
		}
	}

	return query, nil
}

type PaginationQueryImpl struct {
	queries.PaginationQuery
}

func (q *PaginationQueryImpl) GetSpec(page, limit int) repo.QuerySpec {
	return &PaginationSpec{page: page, limit: limit}
}
