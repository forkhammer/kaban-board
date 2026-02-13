package queries

import "main/internal/domain/repo"

type CommonQuery interface {
	PaginationSpec(page, limit int) repo.QuerySpec
	OrderSpec(order string) repo.QuerySpec
}
