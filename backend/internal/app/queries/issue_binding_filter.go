package queries

import "main/internal/domain/repo"

type IssueBindingFilter struct {
	Issues     []uint
	AssigneeId *uint
	TeamId     *uint
	GroupId    *uint
	SprintId   *uint
	ProjectId  *uint
	Search     *string
}

type IssueBindingQuery interface {
	GetSpec(filter IssueBindingFilter) repo.QuerySpec
}
