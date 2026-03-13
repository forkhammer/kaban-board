package repo

import (
	"errors"
	"strconv"
	"time"

	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/gorm"
)

type SprintQuarterRow struct {
	StartDate time.Time
}

type SprintRepository struct {
	conn     interfaces.ConnectionInterface `di.inject:"db"`
	teamRepo *TeamRepository                `di.inject:"TeamRepository"`
}

func (r *SprintRepository) Get(id domain.SprintId) (*domain.Sprint, error) {
	sprint := &models.Sprint{}
	if err := r.getQuery().Where("sprints.id = ?", id).First(sprint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_pkg.NewNotFoundError("Sprint", strconv.Itoa(int(id)), err)
		}
		return nil, err
	}
	return r.toDomainSprint(sprint)
}

func (r *SprintRepository) List(spec repo.QuerySpec) ([]domain.Sprint, error) {
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

	return domainSprints, nil
}

func (r *SprintRepository) Count(spec repo.QuerySpec) (int, error) {
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
		Id:            domain.SprintId(sprint.Id),
		Status:        domain.SprintStatus(sprint.Status),
		Title:         sprint.Title,
		StartDate:     sprint.StartDate,
		EndDate:       sprint.EndDate,
		HoursPerUser:  sprint.HoursPerUser,
		Team:          *team,
		CountBindings: sprint.CountBindings,
	}, nil
}

func (r *SprintRepository) toSprint(sprint *domain.Sprint) *models.Sprint {
	return &models.Sprint{
		Id:           uint(sprint.Id),
		Status:       string(sprint.Status),
		Title:        sprint.Title,
		StartDate:    sprint.StartDate,
		EndDate:      sprint.EndDate,
		HoursPerUser: sprint.HoursPerUser,
		TeamId:       uint(sprint.Team.Id),
	}
}

func (r *SprintRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Sprint{}).
		Select("sprints.*", "(SELECT COUNT(*) FROM issue_bindings WHERE sprint_id = sprints.id) AS count_bindings").
		Joins("Team")
}

func (r *SprintRepository) GetQuarters() ([]domain.Quarter, error) {
	rows := make([]SprintQuarterRow, 0)
	err := r.conn.GetEngine().Raw("SELECT DISTINCT start_date FROM sprints ORDER BY start_date").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	quarters := utils.Map(rows, func(row SprintQuarterRow) domain.Quarter {
		return *domain.NewQuarterFromDate(row.StartDate)
	})

	quarters = utils.Filter(quarters, func(quarter domain.Quarter) bool {
		return quarter.Validate() == nil
	})

	quarters = utils.Unique(quarters, func(quarter domain.Quarter) string {
		return quarter.GetId()
	})
	return quarters, nil
}
