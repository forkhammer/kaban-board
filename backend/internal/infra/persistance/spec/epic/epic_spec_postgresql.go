package epic_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type EpicFilterSpecPostgresql struct {
	repo.QuerySpec
	Filter queries.EpicFilter
}

func (s *EpicFilterSpecPostgresql) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.ProjectId != nil {
		query = query.Where("project_id = ?", s.Filter.ProjectId)
	}

	if s.Filter.Search != nil {
		query = query.Where(
			"(epics.title ILIKE ?) OR (epics.iid = ?)",
			fmt.Sprintf("%%%s%%", *s.Filter.Search),
			*s.Filter.Search,
		)
	}

	return query, nil
}

type EpicQueryImplPostgresql struct {
	queries.EpicQuery
}

func (b *EpicQueryImplPostgresql) GetSpec(filter queries.EpicFilter) repo.QuerySpec {
	return &EpicFilterSpecPostgresql{Filter: filter}
}
