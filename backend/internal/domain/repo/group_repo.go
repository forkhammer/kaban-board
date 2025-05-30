package repo

import (
	"main/internal/domain/models"
)

type GroupRepo interface {
	RWRepo[models.Group, models.GroupId]
}
