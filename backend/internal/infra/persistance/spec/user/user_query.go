package user_spec

import (
	"main/internal/app/queries"
	"main/internal/domain/repo"

	"gorm.io/gorm"
)

type OnlyVisibleUserSpec struct {
	repo.QuerySpec
}

func (s *OnlyVisibleUserSpec) Apply(conn any) (any, error) {
	query := conn.(*gorm.DB)
	return query.Where("is_visible = ?", true), nil
}

type UserQueryImpl struct {
	queries.UserQuery
}

func (q *UserQueryImpl) OnlyVisible() repo.QuerySpec {
	return &OnlyVisibleUserSpec{}
}
