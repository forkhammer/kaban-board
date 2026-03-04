package issuebinding_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type IssueBindingFilterSpec struct {
	repo.QuerySpec
	Filter queries.IssueBindingFilter
}

func (s *IssueBindingFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if len(s.Filter.Issues) > 0 {
		query = query.Where("issue_bindings.issue_id IN ?", s.Filter.Issues)
	}

	if s.Filter.SprintId != nil {
		query = query.Where("issue_bindings.sprint_id = ?", s.Filter.SprintId)
	}

	if s.Filter.AssigneeId != nil {
		query = query.
			Where("issue_bindings.assignee_id = @user_id", map[string]any{
				"user_id": s.Filter.AssigneeId,
			})
	}

	if s.Filter.GroupId != nil {
		query = query.
			Joins("LEFT JOIN user_groups ON user_groups.user_id = issue_bindings.assignee_id OR (user_groups.user_id = assignees.user_id AND issue_bindings.assignee_id IS NULL)").
			Where("user_groups.group_id = ?", s.Filter.GroupId)
	}

	if s.Filter.TeamId != nil {
		query = query.
			Joins("LEFT JOIN issues ON issues.id = issue_bindings.issue_id").
			Joins("LEFT JOIN projects ON issues.project_id = projects.id").
			Where("projects.team_id = ?", s.Filter.TeamId)
	}

	if s.Filter.ProjectId != nil {
		query = query.Where("issues.project_id = ?", s.Filter.ProjectId)
	}
	if s.Filter.Search != nil {
		query = query.Where("(issues.title LIKE ?) OR (issues.iid = ?)", fmt.Sprintf("%%%s%%", *s.Filter.Search), *s.Filter.Search)
	}

	return query, nil
}

type IssueBindingQueryImpl struct {
	queries.IssueBindingQuery
}

func (b *IssueBindingQueryImpl) GetSpec(filter queries.IssueBindingFilter) repo.QuerySpec {
	return &IssueBindingFilterSpec{Filter: filter}
}
