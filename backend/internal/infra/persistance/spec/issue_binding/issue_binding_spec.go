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

	if len(s.Filter.SprintStatuses) > 0 {
		query = query.
			Joins("JOIN sprints ON sprints.id = issue_bindings.sprint_id").
			Where("sprints.status IN ?", s.Filter.SprintStatuses)
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
			Joins("LEFT JOIN user_groups ON user_groups.user_id = issue_bindings.assignee_id").
			Where("user_groups.group_id = ?", s.Filter.GroupId)
	}

	if s.Filter.TeamId != nil && s.Filter.SprintId == nil {
		query = query.
			Joins("LEFT JOIN issues ON issues.id = issue_bindings.issue_id").
			Joins("LEFT JOIN projects ON issues.project_id = projects.id").
			Where("projects.team_id = ?", s.Filter.TeamId)
	}

	if s.Filter.ProjectId != nil {
		query = query.Where("issues.project_id = ?", s.Filter.ProjectId)
	}

	if len(s.Filter.AssigneeGroupIds) > 0 {
		query = query.
			Joins("LEFT JOIN user_groups AS ug_filter ON ug_filter.user_id = issue_bindings.assignee_id").
			Where("ug_filter.group_id IN ?", s.Filter.AssigneeGroupIds)
	}

	if s.Filter.Search != nil {
		query = query.Where(
			"(issues.title LIKE ?) OR (issues.iid = ?) OR (issues.web_url = ?)",
			fmt.Sprintf("%%%s%%", *s.Filter.Search),
			*s.Filter.Search,
			*s.Filter.Search,
		)
	}

	return query, nil
}

type IssueBindingQueryImpl struct {
	queries.IssueBindingQuery
}

func (b *IssueBindingQueryImpl) GetSpec(filter queries.IssueBindingFilter) repo.QuerySpec {
	return &IssueBindingFilterSpec{Filter: filter}
}
