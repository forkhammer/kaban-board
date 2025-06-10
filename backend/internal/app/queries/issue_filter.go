package queries

import (
	"main/internal/domain/repo"
)

type IssueFilter struct {
	AssigneeId *uint
	TeamId     *uint
	GroupId    *uint
	SprintId   *uint
	ProjectId  *uint
	Search     *string
}

type IssueQuery interface {
	GetSpec(filter IssueFilter) repo.QuerySpec
}
