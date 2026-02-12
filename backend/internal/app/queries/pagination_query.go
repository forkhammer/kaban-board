package queries

import "main/internal/domain/repo"

type PaginationQuery interface {
	GetSpec(page, limit int) repo.QuerySpec
}
