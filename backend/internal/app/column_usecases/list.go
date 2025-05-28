package column_usecases

import (
	"fmt"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ListColumnsUseCase struct {
	columnRepo repo.ColumnRepo `di.inject:"ColumnRepository"`
}

func (uc *ListColumnsUseCase) Execute() (*[]domain.Column, error) {
	columns, err := uc.columnRepo.List(nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving columns: %w", err)
	}
	return columns, nil
}
