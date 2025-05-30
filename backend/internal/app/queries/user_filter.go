package queries

import "main/internal/domain/repo"

type UserQuery interface {
	OnlyVisible() repo.QuerySpec
}
