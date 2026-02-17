package user_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type OnlyVisibleUserSpec struct {
	repo.QuerySpec
}

func (s *OnlyVisibleUserSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)
	return query.Where("is_visible = ?", true), nil
}

type UserFilterSpec struct {
	repo.QuerySpec
	Filter queries.UserFilter
}

func (s *UserFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.Search != nil {
		query = query.Where("name like ?", fmt.Sprintf("%%%s%%", *s.Filter.Search))
	}

	if s.Filter.TeamId != nil {
		query = query.Where(
			`
				id IN (
					SELECT user_id 
					FROM assignees 
					INNER JOIN issues ON issues.id = assignees.issue_id
					INNER JOIN projects ON projects.id = issues.project_id
					WHERE projects.team_id = ?
				)
			`,
			*s.Filter.TeamId,
		)
	}

	return query, nil
}

type UserQueryImpl struct {
	queries.UserQuery
}

func (q *UserQueryImpl) OnlyVisible() repo.QuerySpec {
	return &OnlyVisibleUserSpec{}
}

func (q *UserQueryImpl) GetSpec(filter queries.UserFilter) repo.QuerySpec {
	return &UserFilterSpec{Filter: filter}
}
