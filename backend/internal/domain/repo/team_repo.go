package repo

import (
	"main/internal/domain/models"
)

type TeamRepo interface {
	RWRepo[models.Team, models.TeamId]
}
