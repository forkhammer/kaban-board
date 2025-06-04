package sprint_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type SprintFilterSpec struct {
	repo.QuerySpec
	Filter queries.SprintFilter
}

func (s *SprintFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.TeamID != nil {
		query = query.Where("team_id = ?", s.Filter.TeamID)
	}

	return query, nil
}

type SprintQueryImpl struct {
	queries.SprintQuery
}

func (b *SprintQueryImpl) GetSpec(filter queries.SprintFilter) repo.QuerySpec {
	return &SprintFilterSpec{Filter: filter}
}
