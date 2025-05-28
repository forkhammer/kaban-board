package repo

import (
	"main/internal/domain/models"
)

type ProjectRepo interface {
	RWRepo[models.Project, models.ProjectId]
}
