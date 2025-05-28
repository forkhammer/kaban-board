package column_usecases

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type RetrieveColumnUseCase struct {
	columnRepo repo.ColumnRepo `di.inject:"ColumnRepository"`
}

func (uc *RetrieveColumnUseCase) Execute(id uint) (*domain.Column, error) {
	return uc.columnRepo.Get(domain.ColumnId(id))
}
