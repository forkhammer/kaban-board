package repo

import (
	"errors"
	domain "main/internal/domain"
	domain_models "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"strconv"

	"gorm.io/gorm"
)

type EpicRepository struct {
	conn        interfaces.ConnectionInterface `di.inject:"db"`
	projectRepo *ProjectRepository             `di.inject:"ProjectRepository"`
}

func (r *EpicRepository) Get(id domain_models.EpicId) (*domain_models.Epic, error) {
	epic := &models.Epic{}
	if err := r.getQuery().Where("epics.id = ?", id).First(epic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("Epic", strconv.Itoa(int(id)), err)
		}
		return nil, err
	}
	return r.toDomainEpic(epic)
}

func (r *EpicRepository) List(spec repo.QuerySpec) ([]domain_models.Epic, error) {
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

	domainEpics := make([]domain_models.Epic, len(epics))
	for i, release := range epics {
		val, err := r.toDomainEpic(&release)
		if err != nil {
			return nil, err
		}
		domainEpics[i] = *val
	}

	return domainEpics, nil
}

func (r *EpicRepository) Count(spec repo.QuerySpec) (int, error) {
	var count int64
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return 0, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *EpicRepository) Create(epic *domain_models.Epic) (*domain_models.Epic, error) {
	model := r.toEpic(epic)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain_models.EpicId(model.Id))
}

func (r *EpicRepository) Update(epic *domain_models.Epic) (*domain_models.Epic, error) {
	model := r.toEpic(epic)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain_models.EpicId(model.Id))
}

func (r *EpicRepository) Delete(id domain_models.EpicId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Epic{}).Error
}

func (r *EpicRepository) GetByExternalId(externalId domain_models.IssueExternalId) (*domain_models.Epic, error) {
	epic := &models.Epic{}
	if err := r.getQuery().Where("epics.external_id = ?", externalId).First(epic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("Epic", string(externalId), err)

		}
		return nil, err
	}
	return r.toDomainEpic(epic)
}

func (r *EpicRepository) toDomainEpic(release *models.Epic) (*domain_models.Epic, error) {
	project, err := r.projectRepo.toDomainProject(&release.Project)
	if err != nil {
		return nil, err
	}

	return &domain_models.Epic{
		Id:         domain_models.EpicId(release.Id),
		ExternalId: domain_models.IssueExternalId(release.ExternalId),
		Title:      release.Title,
		Project:    *project,
	}, nil
}

func (r *EpicRepository) toEpic(epic *domain_models.Epic) *models.Epic {
	return &models.Epic{
		Id:         uint(epic.Id),
		ExternalId: string(epic.ExternalId),
		Title:      epic.Title,
		ProjectId:  uint(epic.Project.Id),
	}
}

func (r *EpicRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Epic{}).Joins("Project")
}
