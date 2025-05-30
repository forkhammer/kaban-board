package models

import "errors"

type ColumnId uint

type Column struct {
	Id     ColumnId
	Name   string
	Labels []LabelId
	Team   *Team
	Order  *int
}

func NewColumn(id ColumnId, name string, labels []LabelId, team *Team, order *int) (*Column, error) {
	column := &Column{
		Id:     id,
		Name:   name,
		Labels: labels,
		Team:   team,
		Order:  order,
	}

	return column, column.Validate()
}

func (c *Column) Validate() error {
	if c.Name == "" {
		return errors.New("Название колонки не может быть пустым")
	}

	return nil
}
