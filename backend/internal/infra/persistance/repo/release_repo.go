package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type ReleaseRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *ReleaseRepository) Get(id domain.ReleaseId) (*domain.Release, error) {
	release := &models.Release{}
	if err := r.getQuery().Where("id = ?", id).First(release).Error; err != nil {
		return nil, err
	}
	return r.toDomainRelease(release), nil
}

func (r *ReleaseRepository) List(spec repo.QuerySpec) (*[]domain.Release, error) {
	releases := make([]models.Release, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&releases).Error; err != nil {
		return nil, err
	}

	domainReleases := make([]domain.Release, len(releases))
	for i, release := range releases {
		domainReleases[i] = *r.toDomainRelease(&release)
	}

	return &domainReleases, nil
}

func (r *ReleaseRepository) Create(release *domain.Release) error {
	model := r.toRelease(release)
	return r.conn.GetEngine().Create(model).Error
}

func (r *ReleaseRepository) Update(release *domain.Release) error {
	model := r.toRelease(release)
	return r.conn.GetEngine().Save(model).Error
}

func (r *ReleaseRepository) Delete(id domain.ReleaseId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Release{}).Error
}

func (r *ReleaseRepository) toDomainRelease(release *models.Release) *domain.Release {
	return &domain.Release{
		Id:        domain.ReleaseId(release.Id),
		Iid:       domain.ReleaseIid(release.Iid),
		Title:     release.Title,
		ProjectId: domain.ProjectId(release.ProjectId),
	}
}

func (r *ReleaseRepository) toRelease(release *domain.Release) *models.Release {
	return &models.Release{
		Id:        string(release.Id),
		Iid:       string(release.Iid),
		Title:     release.Title,
		ProjectId: uint(release.ProjectId),
	}
}

func (r *ReleaseRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Release{}).Preload("Project")
}
