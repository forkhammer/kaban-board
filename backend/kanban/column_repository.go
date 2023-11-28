package kanban

import "main/db"

type ColumnRepository struct{}

func (r *ColumnRepository) GetColumns() ([]Column, error) {
	var columns []Column
	result := db.DefaultConnection.Db.Order("\"order\"").Find(&columns)
	return columns, result.Error
}

func (r *ColumnRepository) GetColumnById(id int) (*Column, error) {
	var column Column
	if result := db.DefaultConnection.Db.Where(&Column{Id: id}).First(&column); result.Error != nil {
		return nil, result.Error
	} else {
		return &column, nil
	}
}

func (r *ColumnRepository) SaveColumn(column *Column) error {
	result := db.DefaultConnection.Db.Save(column)
	return result.Error
}

func (r *ColumnRepository) CreateColumn(column *Column) error {
	result := db.DefaultConnection.Db.Create(column)
	return result.Error
}

func (r *ColumnRepository) DeleteColumn(id int) error {
	result := db.DefaultConnection.Db.Delete(&Column{Id: id})
	return result.Error
}
