package sprint_user_settings_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type SprintUserSettingsFilterSpec struct {
	repo.QuerySpec
	Filter queries.SprintUserSettingsFilter
}

func (s *SprintUserSettingsFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.SprintId != nil {
		query = query.Where("sprint_id = ?", s.Filter.SprintId)
	}

	if s.Filter.UserId != nil {
		query = query.Where("user_id = ?", s.Filter.UserId)
	}

	return query, nil
}

type SprintUserSettingsQueryImpl struct {
	queries.SprintUserSettingsQuery
}

func (b *SprintUserSettingsQueryImpl) GetSpec(filter queries.SprintUserSettingsFilter) repo.QuerySpec {
	return &SprintUserSettingsFilterSpec{Filter: filter}
}
