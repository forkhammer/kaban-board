package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type EpicUseCases struct {
	epicRepo  repo.EpicRepo     `di.inject:"EpicRepository"`
	epicQuery queries.EpicQuery `di.inject:"EpicQuery"`
}

func (uc *EpicUseCases) GetEpics(filter *queries.EpicFilter) (*[]domain.Epic, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.epicQuery.GetSpec(*filter)
	}

	return uc.epicRepo.List(spec)
}

func (uc *EpicUseCases) GetEpic(id uint) (*domain.Epic, error) {
	return uc.epicRepo.Get(domain.EpicId(id))
}
