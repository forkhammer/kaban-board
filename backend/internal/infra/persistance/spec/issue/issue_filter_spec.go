package issue_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type IssueFilterSpec struct {
	repo.QuerySpec
	Filter queries.IssueFilter
}

func (s *IssueFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.AssigneeId != nil {
		query = query.
			Joins("left join assignees on issues.id = assignees.issue_id").
			Where("assignees.user_id = ?", s.Filter.AssigneeId)
	}
	if s.Filter.GroupId != nil {
		query = query.
			Joins("left join assignees on issues.id = assignees.issue_id").
			Joins("left join users on assignees.user_id = users.id").
			Joins("left join user_groups on users.id = user_groups.user_id").
			Where("user_groups.group_id = ?", s.Filter.GroupId)
	}
	if s.Filter.TeamId != nil {
		query = query.
			Joins("left join projects on issues.project_id = projects.id").
			Where("projects.team_id = ?", s.Filter.TeamId)
	}
	if s.Filter.SprintId != nil {

	}

	return query, nil
}

type IssueQueryImpl struct {
	queries.IssueQuery
}

func (b *IssueQueryImpl) GetSpec(filter queries.IssueFilter) repo.QuerySpec {
	return &IssueFilterSpec{Filter: filter}
}
