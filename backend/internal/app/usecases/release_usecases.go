package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ReleaseUseCases struct {
	releaseRepo  repo.ReleaseRepo     `di.inject:"ReleaseRepository"`
	releaseQuery queries.ReleaseQuery `di.inject:"ReleaseQuery"`
}

func (uc *ReleaseUseCases) GetReleases(filter *queries.ReleaseFilter) ([]domain.Release, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.releaseQuery.GetSpec(*filter)
	}

	return uc.releaseRepo.List(spec)
}

func (uc *ReleaseUseCases) GetRelease(id uint) (*domain.Release, error) {
	return uc.releaseRepo.Get(domain.ReleaseId(id))
}
