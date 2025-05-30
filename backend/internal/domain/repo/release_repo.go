package repo

import (
	"main/internal/domain/models"
)

type ReleaseRepo interface {
	RWRepo[models.Release, models.ReleaseId]
}
