package queries

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ColumnFilter struct {
	Ids []domain.ColumnId
}

type ColumnQuery interface {
	GetSpec(filter ColumnFilter) repo.QuerySpec
}
