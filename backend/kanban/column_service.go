package kanban

import (
	"errors"
)

type ColumnService struct {
	columnRepository ColumnRepository
}

func (s *ColumnService) GetAllColumns() ([]Column, error) {
	return s.columnRepository.GetColumns()
}

func (s *ColumnService) GetColumnById(id int) (*Column, error) {
	return s.columnRepository.GetColumnById(id)
}

func (s *ColumnService) UpdateColumn(id int, data *UpdateColumnRequest) (*Column, error) {
	column, err := s.columnRepository.GetColumnById(id)

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
		return column, nil
	}
}

func (s *ColumnService) CreateColumn(data *CreateColumnRequest) (*Column, error) {
	if data.Name == "" {
		return nil, errors.New("Название колонки не может быть пустым")
	}

	column := Column{
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
	return s.columnRepository.DeleteColumn(id)
}
