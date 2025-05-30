package repo

import (
	"main/internal/domain/models"
)

type LabelRepo interface {
	RWRepo[models.Label, models.LabelId]
}
