package queries

import "main/internal/domain/repo"

type SprintUserSettingsFilter struct {
	SprintId *uint
	UserId   *uint
}

type SprintUserSettingsQuery interface {
	GetSpec(filter SprintUserSettingsFilter) repo.QuerySpec
}
