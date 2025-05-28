package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ColumnRepository struct {
	conn     interfaces.ConnectionInterface `di.inject:"db"`
	teamRepo *TeamRepository                `di.inject:"TeamRepository"`
}

func (r *ColumnRepository) Get(id domain.ColumnId) (*domain.Column, error) {
	сolumn := &models.Column{}
	if err := r.getQuery().Where("id = ?", id).First(сolumn).Error; err != nil {
		return nil, err
	}
	model, err := r.toDomainColumn(сolumn)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *ColumnRepository) List(spec repo.QuerySpec) (*[]domain.Column, error) {
	сolumns := make([]models.Column, 0)
	query := r.getQuery().Model(&models.Column{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&сolumns).Error; err != nil {
		return nil, err
	}

	domainColumns := make([]domain.Column, len(сolumns))
	for i, Column := range сolumns {
		if elem, err := r.toDomainColumn(&Column); err != nil {
			return nil, err
		} else {
			domainColumns[i] = *elem
		}
	}

	return &domainColumns, nil
}

func (r *ColumnRepository) Create(column *domain.Column) (*domain.Column, error) {
	model, err := r.toColumn(column)
	if err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Create(model).Error; err != nil {
		return nil, err
	}

	result, err := r.toDomainColumn(model)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ColumnRepository) Update(column *domain.Column) (*domain.Column, error) {
	model, err := r.toColumn(column)
	if err != nil {
		return nil, err
	}
	if err := r.conn.GetEngine().Save(model).Error; err != nil {
		return nil, err
	}

	result, err := r.toDomainColumn(model)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ColumnRepository) Delete(id domain.ColumnId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Column{}).Error
}

func (r *ColumnRepository) toDomainColumn(column *models.Column) (*domain.Column, error) {
	var team *domain.Team
	if column.Team != nil {
		team = r.teamRepo.ToDomainTeam(column.Team)
	}

	result, err := domain.NewColumn(
		domain.ColumnId(column.Id),
		column.Name,
		utils.Map(column.Labels, func(labelId string) domain.LabelId {
			return domain.LabelId(labelId)
		}),
		team,
		column.Order,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ColumnRepository) toColumn(column *domain.Column) (*models.Column, error) {
	var teamId *uint
	if column.Team != nil {
		val := uint(column.Team.Id)
		teamId = &val
	}

	return &models.Column{
		Id:     int(column.Id),
		Name:   column.Name,
		TeamId: teamId,
		Labels: datatypes.NewJSONSlice(utils.Map(column.Labels, func(labelId domain.LabelId) string {
			return string(labelId)
		})),
		Order: column.Order,
	}, nil
}

func (r *ColumnRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Column{}).Preload("Team")
}
