package user_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"
	"strings"

	"gorm.io/gorm"
)

type UserFilterSpecMysql struct {
	repo.QuerySpec
	Filter queries.UserFilter
}

func (s *UserFilterSpecMysql) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.Search != nil {
		searchText := fmt.Sprintf("%%%s%%", strings.ToLower(*s.Filter.Search))
		query = query.Where("(lower(users.name) like ?) or (lower(users.username) like ?)", searchText, searchText)
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

type UserQueryImplMysql struct {
	queries.UserQuery
}

func (q *UserQueryImplMysql) OnlyVisible() repo.QuerySpec {
	return &OnlyVisibleUserSpec{}
}

func (q *UserQueryImplMysql) GetSpec(filter queries.UserFilter) repo.QuerySpec {
	return &UserFilterSpecMysql{Filter: filter}
}
