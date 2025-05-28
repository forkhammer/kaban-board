package repo

import (
	"main/internal/domain/models"
)

type UserRepo interface {
	RWRepo[models.User, models.UserId]
}
