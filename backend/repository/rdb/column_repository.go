package rdb

import (
	"main/tools"

	"gorm.io/gorm"
)

type ColumnRepository struct {
	connection tools.ConnectionInterface
}

func (r *ColumnRepository) GetColumns(to interface{}) error {
	result := r.getDb().Order("\"order\"").Find(to)
	return result.Error
}

func (r *ColumnRepository) GetColumnById(to interface{}, id int) error {
	if result := r.getDb().Where("id = ?", id).First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *ColumnRepository) SaveColumn(column interface{}) error {
	result := r.getDb().Save(column)
	return result.Error
}

func (r *ColumnRepository) CreateColumn(column interface{}) error {
	result := r.getDb().Create(column)
	return result.Error
}

func (r *ColumnRepository) DeleteColumn(id int) error {
	result := r.getDb().Delete("id = ?", id)
	return result.Error
}

func (r *ColumnRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
