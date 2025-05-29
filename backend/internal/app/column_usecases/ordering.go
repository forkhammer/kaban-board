package column_usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type ColumnOrdering struct {
	Id    uint
	Order int
}

type ColumnOrderingSet = []ColumnOrdering

type OrderingColumnUseCase struct {
	columnRepo  repo.ColumnRepo     `di.inject:"ColumnRepository"`
	columnQuery queries.ColumnQuery `di.inject:"ColumnQuery"`
}

func (uc *OrderingColumnUseCase) Execute(ordering ColumnOrderingSet) (*[]domain.Column, error) {
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
