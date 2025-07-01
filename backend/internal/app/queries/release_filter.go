package queries

import "main/internal/domain/repo"

type ReleaseFilter struct {
	Search    *string
	ProjectId *uint
}

type ReleaseQuery interface {
	GetSpec(filter ReleaseFilter) repo.QuerySpec
}
