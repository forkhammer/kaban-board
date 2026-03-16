package project_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type ProjectFilterSpecPostgresql struct {
	repo.QuerySpec
	Filter queries.ProjectFilter
}

func (s *ProjectFilterSpecPostgresql) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.TeamID != nil {
		query = query.Where("team_id = ?", s.Filter.TeamID)
	}

	if s.Filter.Search != nil {
		searchText := *s.Filter.Search
		query = query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", searchText))
	}

	return query, nil
}

type ProjectQueryImplPostgresql struct {
	queries.ProjectQuery
}

func (b *ProjectQueryImplPostgresql) GetSpec(filter queries.ProjectFilter) repo.QuerySpec {
	return &ProjectFilterSpecPostgresql{Filter: filter}
}
