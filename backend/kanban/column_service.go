package kanban

import (
	"errors"
	"main/repository"
	"main/repository/models"
	"main/tools"
)

type ColumnService struct {
	columnRepository repository.ColumnRepositoryInterface `di.inject:"columnRepository"`
}

func (s *ColumnService) GetAllColumns() ([]models.Column, error) {
	var columns []models.Column
	err := s.columnRepository.GetColumns(&columns)
	return columns, err
}

func (s *ColumnService) GetColumnById(id int) (*models.Column, error) {
	var column models.Column
	err := s.columnRepository.GetColumnById(&column, id)
	return &column, err
}

func (s *ColumnService) UpdateColumn(id int, data *UpdateColumnRequest) (*models.Column, error) {
	var column models.Column
	err := s.columnRepository.GetColumnById(&column, id)

	if err != nil {
		return nil, err
	}

	if data.Name == "" {
		return nil, errors.New("Название колонки не может быть пустым")
	}

	column.Name = data.Name
	column.Labels = data.Labels
	column.TeamId = data.TeamId
	if err = s.columnRepository.SaveColumn(column); err != nil {
		return nil, err
	} else {
		return &column, nil
	}
}

func (s *ColumnService) CreateColumn(data *CreateColumnRequest) (*models.Column, error) {
	if data.Name == "" {
		return nil, errors.New("Название колонки не может быть пустым")
	}

	column := models.Column{
		Name:   data.Name,
		Labels: data.Labels,
		TeamId: data.TeamId,
	}

	if err := s.columnRepository.CreateColumn(&column); err != nil {
		return nil, err
	} else {
		return &column, nil
	}
}

func (s *ColumnService) DeleteColumnById(id int) error {
	return s.columnRepository.DeleteColumn(&models.Column{Id: id})
}

func (s *ColumnService) SaveOrdering(request []SetColumnOrderRequest) ([]models.Column, error) {
	columns, err := s.GetAllColumns()
	result := make([]models.Column, 0)

	if err != nil {
		return nil, err
	}

	for _, req := range request {
		column := tools.Find[models.Column](columns, func(column models.Column) bool {
			return column.Id == req.Id
		})
		if column != nil {
			order := req.Order
			column.Order = &order
			if err = s.columnRepository.SaveColumn(column); err != nil {
				return nil, err
			}
			result = append(result, *column)
		}
	}

	return result, nil
}
