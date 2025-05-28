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
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *ColumnRepository) Get(id domain.ColumnId) (*domain.Column, error) {
	сolumn := &models.Column{}
	if err := r.conn.GetEngine().Where("id = ?", id).First(сolumn).Error; err != nil {
		return nil, err
	}
	model, err := r.toDomainColumn(сolumn)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *ColumnRepository) GetByUsername(username string) (*domain.Column, error) {
	сolumn := &models.Column{}
	if err := r.conn.GetEngine().Where("username = ?", username).First(сolumn).Error; err != nil {
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
	query := r.conn.GetEngine().Model(&models.Column{})

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

func (r *ColumnRepository) Create(Column *domain.Column) error {
	if model, err := r.toColumn(Column); err != nil {
		return err
	} else {
		return r.conn.GetEngine().Create(model).Error
	}
}

func (r *ColumnRepository) Update(Column *domain.Column) error {
	if model, err := r.toColumn(Column); err != nil {
		return err
	} else {
		return r.conn.GetEngine().Save(model).Error
	}
}

func (r *ColumnRepository) Delete(id domain.ColumnId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Column{}).Error
}

func (r *ColumnRepository) toDomainColumn(column *models.Column) (*domain.Column, error) {
	teamId, err := utils.IntToUintPtr(column.TeamId)
	if err != nil {
		return nil, err
	}

	return &domain.Column{
		Id:     domain.ColumnId(column.Id),
		Name:   column.Name,
		TeamId: (*domain.TeamId)(teamId),
		Labels: utils.Map(column.Labels, func(labelId string) domain.LabelId {
			return domain.LabelId(labelId)
		}),
		Order: column.Order,
	}, nil
}

func (r *ColumnRepository) toColumn(column *domain.Column) (*models.Column, error) {
	teamId, err := utils.UintToIntPtr((*uint)(column.TeamId))
	if err != nil {
		return nil, nil
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
