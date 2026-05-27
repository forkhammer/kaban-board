package queries

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"
)

type IssueBindingFilter struct {
	Issues           []uint
	SprintStatuses   []models.SprintStatus
	AssigneeId       *uint
	TeamId           *uint
	GroupId          *uint
	SprintId         *uint
	ProjectId        *uint
	Search           *string
	AssigneeGroupIds []uint
}

type IssueBindingQuery interface {
	GetSpec(filter IssueBindingFilter) repo.QuerySpec
}
