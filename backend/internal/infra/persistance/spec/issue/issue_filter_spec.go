package issue_spec

import (
	"fmt"
	"main/internal/app/queries"
	"main/internal/domain/repo"
	"strings"

	"gorm.io/gorm"
)

type IssueFilterSpec struct {
	repo.QuerySpec
	Filter queries.IssueFilter
}

func (s *IssueFilterSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)

	if s.Filter.AssigneeId != nil && s.Filter.SprintId == nil {
		query = query.
			Joins("left join assignees on issues.id = assignees.issue_id").
			Where("assignees.user_id = ?", s.Filter.AssigneeId)
	}
	if s.Filter.GroupId != nil && s.Filter.SprintId == nil {
		query = query.
			Joins("left join assignees on issues.id = assignees.issue_id").
			Joins("left join users on assignees.user_id = users.id").
			Joins("left join user_groups on users.id = user_groups.user_id").
			Where("user_groups.group_id = ?", s.Filter.GroupId)
	}
	if s.Filter.TeamId != nil && s.Filter.SprintId == nil {
		query = query.
			Joins("left join projects on issues.project_id = projects.id").
			Where("projects.team_id = ?", s.Filter.TeamId)
	}
	if s.Filter.SprintId != nil {
		query = query.
			Select("issues.*, issue_bindings.id as binding_id").
			Joins("left join issue_bindings on issues.id = issue_bindings.issue_id").
			Where("issue_bindings.sprint_id = ?", s.Filter.SprintId)

		if s.Filter.AssigneeId != nil {
			query = query.
				Joins("left join assignees on issues.id = assignees.issue_id").
				Where("issue_bindings.assignee_id = @user_id OR (assignees.user_id = @user_id AND issue_bindings.assignee_id IS NULL)", map[string]any{
					"user_id": s.Filter.AssigneeId,
				})
		}

		if s.Filter.GroupId != nil {
			query = query.
				Joins("left join assignees on issues.id = assignees.issue_id").
				Joins("left join users on assignees.user_id = users.id").
				Joins("left join user_groups on user_groups.user_id = issue_bindings.assignee_id OR (user_groups.user_id = users.id AND issue_bindings.assignee_id IS NULL)").
				Where("user_groups.group_id = ?", s.Filter.GroupId)
		}
	}
	if s.Filter.ProjectId != nil {
		query = query.Where("issues.project_id = ?", s.Filter.ProjectId)
	}
	if s.Filter.Search != nil {
		searchText := strings.ToLower(*s.Filter.Search)
		query = query.Where("(lower_unicode(issues.title) like ?) or (issues.iid = ?)", fmt.Sprintf("%%%s%%", searchText), searchText)
	}

	return query, nil
}

type IssueQueryImpl struct {
	queries.IssueQuery
}

func (b *IssueQueryImpl) GetSpec(filter queries.IssueFilter) repo.QuerySpec {
	return &IssueFilterSpec{Filter: filter}
}
