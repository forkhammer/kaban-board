package usecases

import (
	"fmt"
	"main/internal/app/queries"
	app_services "main/internal/app/services"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type CreateColumnRequest struct {
	Name   string
	Labels []string
	TeamId *uint
}

type UpdateColumnRequest struct {
	Id     uint
	Name   string
	Labels []string
	TeamId *uint
}

type ColumnOrdering struct {
	Id    uint
	Order int
}

type ColumnOrderingSet = []ColumnOrdering

type ColumnUseCases struct {
	columnRepo   repo.ColumnRepo            `di.inject:"ColumnRepository"`
	teamRepo     repo.TeamRepo              `di.inject:"TeamRepository"`
	labelRepo    repo.LabelRepo             `di.inject:"LabelRepository"`
	labelService *app_services.LabelService `di.inject:"LabelService"`
	columnQuery  queries.ColumnQuery        `di.inject:"ColumnQuery"`
}

func (uc *ColumnUseCases) Create(request *CreateColumnRequest) (*domain.Column, error) {
	var team *domain.Team
	var err error
	if request.TeamId != nil {
		team, err = uc.teamRepo.Get(domain.TeamId(*request.TeamId))
		if err != nil {
			return nil, err
		}
	}
	order := 0

	existLabels, diff, err := uc.labelService.ExistNames(request.Labels)
	if err != nil {
		return nil, err
	}
	if !existLabels {
		return nil, fmt.Errorf("Labels not found: %v", diff)
	}

	column, err := domain.NewColumn(
		0,
		request.Name,
		utils.Map(request.Labels, func(label string) domain.LabelId {
			return domain.LabelId(label)
		}),
		team,
		&order,
	)

	if err := column.Validate(); err != nil {
		return nil, err
	}

	column, err = uc.columnRepo.Create(column)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (c *ColumnUseCases) Delete(id uint) error {
	_, err := c.columnRepo.Get(domain.ColumnId(id))
	if err != nil {
		return err
	}

	return c.columnRepo.Delete(domain.ColumnId(id))
}

func (uc *ColumnUseCases) List() (*[]domain.Column, error) {
	columns, err := uc.columnRepo.List(nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving columns: %w", err)
	}
	return columns, nil
}

func (uc *ColumnUseCases) Ordering(ordering ColumnOrderingSet) (*[]domain.Column, error) {
	columnIds := utils.Map(ordering, func(o ColumnOrdering) domain.ColumnId {
		return domain.ColumnId(o.Id)
	})
	columns, err := uc.columnRepo.List(uc.columnQuery.GetSpec(queries.ColumnFilter{Ids: columnIds}))
	if err != nil {
		return nil, err
	}

	result := make([]domain.Column, len(*columns))

	for index, row := range *columns {
		column := &row
		o := utils.Find(ordering, func(o ColumnOrdering) bool {
			return o.Id == uint(column.Id)
		})
		if o != nil {
			column.Order = &o.Order
		}

		if err := column.Validate(); err != nil {
			return nil, err
		}

		if column, err = uc.columnRepo.Update(column); err != nil {
			return nil, err
		}

		result[index] = *column
	}

	return &result, nil
}

func (uc *ColumnUseCases) Retrieve(id uint) (*domain.Column, error) {
	return uc.columnRepo.Get(domain.ColumnId(id))
}

func (uc *ColumnUseCases) Update(request *UpdateColumnRequest) (*domain.Column, error) {
	existLabels, diff, err := uc.labelService.ExistNames(request.Labels)
	if err != nil {
		return nil, err
	}
	if !existLabels {
		return nil, fmt.Errorf("Labels not found: %v", diff)
	}

	column, err := uc.columnRepo.Get(domain.ColumnId(request.Id))
	if err != nil {
		return nil, err
	}
	column.Name = request.Name
	column.Labels = utils.Map(request.Labels, func(label string) domain.LabelId {
		return domain.LabelId(label)
	})

	var team *domain.Team
	if request.TeamId != nil {
		team, err = uc.teamRepo.Get(domain.TeamId(*request.TeamId))
		if err != nil {
			return nil, err
		}
	}
	column.Team = team

	if err := column.Validate(); err != nil {
		return nil, err
	}

	column, err = uc.columnRepo.Update(column)
	if err != nil {
		return nil, err
	}

	return column, nil
}
