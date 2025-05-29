package queries

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type GroupFilter struct {
	Ids []domain.GroupId
}

type GroupQuery interface {
	GetSpec(filter GroupFilter) repo.QuerySpec
}
