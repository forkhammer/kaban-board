package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type SprintRepository struct {
	conn     interfaces.ConnectionInterface `di.inject:"db"`
	teamRepo *TeamRepository                `di.inject:"TeamRepository"`
}

func (r *SprintRepository) Get(id domain.SprintId) (*domain.Sprint, error) {
	sprint := &models.Sprint{}
	if err := r.getQuery().Where("id = ?", id).First(sprint).Error; err != nil {
		return nil, err
	}
	return r.toDomainSprint(sprint)
}

func (r *SprintRepository) List(spec repo.QuerySpec) (*[]domain.Sprint, error) {
	sprints := make([]models.Sprint, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Order("start_date ASC").Find(&sprints).Error; err != nil {
		return nil, err
	}

	domainSprints := make([]domain.Sprint, len(sprints))
	for i, user := range sprints {
		val, err := r.toDomainSprint(&user)
		if err != nil {
			return nil, err
		}
		domainSprints[i] = *val
	}

	return &domainSprints, nil
}

func (r *SprintRepository) Create(sprint *domain.Sprint) (*domain.Sprint, error) {
	model := r.toSprint(sprint)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}

	return r.Get(domain.SprintId(model.Id))
}

func (r *SprintRepository) Update(sprint *domain.Sprint) (*domain.Sprint, error) {
	model := r.toSprint(sprint)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}

	return r.Get(domain.SprintId(model.Id))
}

func (r *SprintRepository) Delete(id domain.SprintId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Sprint{}).Error
}

func (r *SprintRepository) toDomainSprint(sprint *models.Sprint) (*domain.Sprint, error) {
	team, err := r.teamRepo.ToDomainTeam(&sprint.Team)
	if err != nil {
		return nil, err
	}
	return &domain.Sprint{
		Id:           domain.SprintId(sprint.Id),
		IsCompleted:  sprint.IsCompleted,
		Title:        sprint.Title,
		StartDate:    sprint.StartDate,
		EndDate:      sprint.EndDate,
		HoursPerUser: sprint.HoursPerUser,
		Team:         *team,
	}, nil
}

func (r *SprintRepository) toSprint(sprint *domain.Sprint) *models.Sprint {
	return &models.Sprint{
		Id:           uint(sprint.Id),
		IsCompleted:  sprint.IsCompleted,
		Title:        sprint.Title,
		StartDate:    sprint.StartDate,
		EndDate:      sprint.EndDate,
		HoursPerUser: sprint.HoursPerUser,
		TeamId:       uint(sprint.Team.Id),
	}
}

func (r *SprintRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Sprint{}).Preload("Team")
}
