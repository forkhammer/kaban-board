package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type ReleaseRepository struct {
	conn        interfaces.ConnectionInterface `di.inject:"db"`
	projectRepo *ProjectRepository             `di.inject:"ProjectRepository"`
}

func (r *ReleaseRepository) Get(id domain.ReleaseId) (*domain.Release, error) {
	release := &models.Release{}
	if err := r.getQuery().Where("id = ?", id).First(release).Error; err != nil {
		return nil, err
	}
	return r.toDomainRelease(release)
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
		val, err := r.toDomainRelease(&release)
		if err != nil {
			return nil, err
		}
		domainReleases[i] = *val
	}

	return &domainReleases, nil
}

func (r *ReleaseRepository) Create(release *domain.Release) (*domain.Release, error) {
	model := r.toRelease(release)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.ReleaseId(model.Id))
}

func (r *ReleaseRepository) Update(release *domain.Release) (*domain.Release, error) {
	model := r.toRelease(release)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.ReleaseId(model.Id))
}

func (r *ReleaseRepository) Delete(id domain.ReleaseId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Release{}).Error
}

func (r *ReleaseRepository) toDomainRelease(release *models.Release) (*domain.Release, error) {
	project, err := r.projectRepo.toDomainProject(&release.Project)
	if err != nil {
		return nil, err
	}

	return &domain.Release{
		Id:      domain.ReleaseId(release.Id),
		Iid:     domain.ReleaseIid(release.Iid),
		Title:   release.Title,
		Project: *project,
		WebPath: release.WebPath,
	}, nil
}

func (r *ReleaseRepository) toRelease(release *domain.Release) *models.Release {
	return &models.Release{
		Id:        string(release.Id),
		Iid:       string(release.Iid),
		Title:     release.Title,
		ProjectId: uint(release.Project.Id),
		WebPath:   release.WebPath,
	}
}

func (r *ReleaseRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Release{}).Preload("Project")
}
