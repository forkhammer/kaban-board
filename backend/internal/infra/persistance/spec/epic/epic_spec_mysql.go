package epic_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"
	"strings"

	"gorm.io/gorm"
)

type EpicFilterSpecMysql struct {
	repo.QuerySpec
	Filter queries.EpicFilter
}

func (s *EpicFilterSpecMysql) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.ProjectId != nil {
		query = query.Where("project_id = ?", s.Filter.ProjectId)
	}

	if s.Filter.Search != nil {
		query = query.Where(
			"(lower(epics.title) like ?) OR (epics.iid = ?)",
			fmt.Sprintf("%%%s%%", strings.ToLower(*s.Filter.Search)),
			*s.Filter.Search,
		)
	}

	return query, nil
}

type EpicQueryImplMysql struct {
	queries.EpicQuery
}

func (b *EpicQueryImplMysql) GetSpec(filter queries.EpicFilter) repo.QuerySpec {
	return &EpicFilterSpecMysql{Filter: filter}
}
