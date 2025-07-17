package queries

import "main/internal/domain/repo"

type EpicFilter struct {
	Search    *string
	ProjectId *uint
}

type EpicQuery interface {
	GetSpec(filter EpicFilter) repo.QuerySpec
}
