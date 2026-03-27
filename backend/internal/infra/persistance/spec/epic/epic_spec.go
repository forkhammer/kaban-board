package epic_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"
	"strings"

	"gorm.io/gorm"
)

type EpicFilterSpec struct {
	repo.QuerySpec
	Filter queries.EpicFilter
}

func (s *EpicFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.ProjectId != nil {
		query = query.Where("project_id = ?", s.Filter.ProjectId)
	}

	if s.Filter.Search != nil {
		query = query.Where(
			"(lower_unicode(epics.title) like ?) OR (epics.iid = ?)",
			fmt.Sprintf("%%%s%%", strings.ToLower(*s.Filter.Search)),
			*s.Filter.Search,
		)
	}

	return query, nil
}

type EpicQueryImpl struct {
	queries.EpicQuery
}

func (b *EpicQueryImpl) GetSpec(filter queries.EpicFilter) repo.QuerySpec {
	return &EpicFilterSpec{Filter: filter}
}
