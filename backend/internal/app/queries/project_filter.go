package queries

import "main/internal/domain/repo"

type ProjectFilter struct {
	TeamID *uint
	Search *string
}

type ProjectQuery interface {
	GetSpec(filter ProjectFilter) repo.QuerySpec
}
