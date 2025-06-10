package project_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type ProjectFilterSpec struct {
	repo.QuerySpec
	Filter queries.ProjectFilter
}

func (s *ProjectFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.TeamID != nil {
		query = query.Where("team_id = ?", s.Filter.TeamID)
	}

	if s.Filter.Search != nil {
		query = query.Where("name like ?", fmt.Sprintf("%%%s%%", *s.Filter.Search))
	}

	return query, nil
}

type ProjectQueryImpl struct {
	queries.ProjectQuery
}

func (b *ProjectQueryImpl) GetSpec(filter queries.ProjectFilter) repo.QuerySpec {
	return &ProjectFilterSpec{Filter: filter}
}
