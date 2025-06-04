package queries

import "main/internal/domain/repo"

type SprintFilter struct {
	TeamID    *int
	QuarterId *string
}

type SprintQuery interface {
	GetSpec(filter SprintFilter) repo.QuerySpec
}
