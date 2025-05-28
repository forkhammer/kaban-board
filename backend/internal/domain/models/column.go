package models

type ColumnId uint

type Column struct {
	Id     ColumnId
	Name   string
	Labels []LabelId
	TeamId *TeamId
	Order  *int
}
