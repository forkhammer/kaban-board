package release_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type ReleaseFilterSpec struct {
	repo.QuerySpec
	Filter queries.ReleaseFilter
}

func (s *ReleaseFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.ProjectId != nil {
		query = query.Where("project_id = ?", s.Filter.ProjectId)
	}

	if s.Filter.Search != nil {
		query = query.Where("title like ?", fmt.Sprintf("%%%s%%", *s.Filter.Search))
	}

	return query, nil
}

type ReleaseQueryImpl struct {
	queries.ReleaseQuery
}

func (b *ReleaseQueryImpl) GetSpec(filter queries.ReleaseFilter) repo.QuerySpec {
	return &ReleaseFilterSpec{Filter: filter}
}
