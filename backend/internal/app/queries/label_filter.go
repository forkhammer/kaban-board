package queries

import "main/internal/domain/repo"

type LabelFilter struct {
	Names []string
}

type LabelQuery interface {
	GetSpec(filter LabelFilter) repo.QuerySpec
}
