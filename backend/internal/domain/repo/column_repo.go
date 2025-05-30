package repo

import (
	"main/internal/domain/models"
)

type ColumnRepo interface {
	RWRepo[models.Column, models.ColumnId]
}
