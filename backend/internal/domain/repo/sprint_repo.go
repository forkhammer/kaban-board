package repo

import (
	"main/internal/domain/models"
)

type SprintRepo interface {
	RWRepo[models.Sprint, models.SprintId]
}
