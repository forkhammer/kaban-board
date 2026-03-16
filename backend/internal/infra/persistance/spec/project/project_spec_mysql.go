package project_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"
	"strings"

	"gorm.io/gorm"
)

type ProjectFilterSpecMysql struct {
	repo.QuerySpec
	Filter queries.ProjectFilter
}

func (s *ProjectFilterSpecMysql) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.TeamID != nil {
		query = query.Where("team_id = ?", s.Filter.TeamID)
	}

	if s.Filter.Search != nil {
		searchText := strings.ToLower(*s.Filter.Search)
		query = query.Where("lower(name) like ?", fmt.Sprintf("%%%s%%", searchText))
	}

	return query, nil
}

type ProjectQueryImplMysql struct {
	queries.ProjectQuery
}

func (b *ProjectQueryImplMysql) GetSpec(filter queries.ProjectFilter) repo.QuerySpec {
	return &ProjectFilterSpecMysql{Filter: filter}
}
