package user_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type UserFilterSpecPostgresql struct {
	repo.QuerySpec
	Filter queries.UserFilter
}

func (s *UserFilterSpecPostgresql) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.Search != nil {
		searchText := fmt.Sprintf("%%%s%%", *s.Filter.Search)
		query = query.Where("(users.name ILIKE ?) or (users.username ILIKE ?)", searchText, searchText)
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

type UserQueryImplPostgresql struct {
	queries.UserQuery
}

func (q *UserQueryImplPostgresql) OnlyVisible() repo.QuerySpec {
	return &OnlyVisibleUserSpec{}
}

func (q *UserQueryImplPostgresql) GetSpec(filter queries.UserFilter) repo.QuerySpec {
	return &UserFilterSpecPostgresql{Filter: filter}
}
