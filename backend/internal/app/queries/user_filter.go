package queries

import "main/internal/domain/repo"

type UserFilter struct {
	TeamId *uint
	Search *string
}

type UserQuery interface {
	GetSpec(filter UserFilter) repo.QuerySpec
	OnlyVisible() repo.QuerySpec
}
