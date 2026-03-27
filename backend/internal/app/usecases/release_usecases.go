package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ReleaseUseCases struct {
	releaseRepo  repo.ReleaseRepo     `di.inject:"ReleaseRepository"`
	releaseQuery queries.ReleaseQuery `di.inject:"ReleaseQuery"`
	commonQuery  queries.CommonQuery  `di.inject:"CommonQuery"`
}

func (uc *ReleaseUseCases) GetReleases(filter *queries.ReleaseFilter) ([]domain.Release, error) {
	var spec repo.QuerySpec
	if filter != nil {
		spec = uc.releaseQuery.GetSpec(*filter)
	}
	spec = repo.And(spec, uc.commonQuery.OrderSpec("releases.sort_index DESC NULLS LAST"))
	return uc.releaseRepo.List(spec)
}

func (uc *ReleaseUseCases) GetRelease(id uint) (*domain.Release, error) {
	return uc.releaseRepo.Get(domain.ReleaseId(id))
}
