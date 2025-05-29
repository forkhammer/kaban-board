package column_usecases

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type DeleteColumUseCase struct {
	columnRepo repo.ColumnRepo `di.inject:"ColumnRepository"`
}

func (c *DeleteColumUseCase) Execute(id uint) error {
	_, err := c.columnRepo.Get(domain.ColumnId(id))
	if err != nil {
		return err
	}

	return c.columnRepo.Delete(domain.ColumnId(id))
}
