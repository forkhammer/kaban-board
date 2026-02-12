package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type EpicRepository struct {
	conn        interfaces.ConnectionInterface `di.inject:"db"`
	projectRepo *ProjectRepository             `di.inject:"ProjectRepository"`
}

func (r *EpicRepository) Get(id domain.EpicId) (*domain.Epic, error) {
	epic := &models.Epic{}
	if err := r.getQuery().Where("id = ?", id).First(epic).Error; err != nil {
		return nil, err
	}
	return r.toDomainEpic(epic)
}

func (r *EpicRepository) List(spec repo.QuerySpec) ([]domain.Epic, error) {
	epics := make([]models.Epic, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&epics).Error; err != nil {
		return nil, err
	}

	domainEpics := make([]domain.Epic, len(epics))
	for i, release := range epics {
		val, err := r.toDomainEpic(&release)
		if err != nil {
			return nil, err
		}
		domainEpics[i] = *val
	}

	return domainEpics, nil
}

func (r *EpicRepository) Create(epic *domain.Epic) (*domain.Epic, error) {
	model := r.toEpic(epic)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.EpicId(model.Id))
}

func (r *EpicRepository) Update(epic *domain.Epic) (*domain.Epic, error) {
	model := r.toEpic(epic)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.EpicId(model.Id))
}

func (r *EpicRepository) Delete(id domain.EpicId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Epic{}).Error
}

func (r *EpicRepository) toDomainEpic(release *models.Epic) (*domain.Epic, error) {
	project, err := r.projectRepo.toDomainProject(&release.Project)
	if err != nil {
		return nil, err
	}

	return &domain.Epic{
		Id:      domain.EpicId(release.Id),
		Title:   release.Title,
		Project: *project,
	}, nil
}

func (r *EpicRepository) toEpic(epic *domain.Epic) *models.Epic {
	return &models.Epic{
		Id:        uint(epic.Id),
		Title:     epic.Title,
		ProjectId: uint(epic.Project.Id),
	}
}

func (r *EpicRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Epic{}).Preload("Project")
}
