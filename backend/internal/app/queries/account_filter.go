package queries

import "main/internal/domain/repo"

type AccountFilter struct {
	GitlabIds []uint
	Ids       []uint
}

type AccountQuery interface {
	GetSpec(filter AccountFilter) repo.QuerySpec
}
