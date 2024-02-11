package rdb

import (
	"main/tools"

	"gorm.io/gorm"
)

type LabelRepository struct {
	connection tools.ConnectionInterface
}

func (r *LabelRepository) GetOrCreate(to, query, attrs interface{}) error {
	if result := r.getDb().Where(query).Attrs(attrs).FirstOrCreate(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *LabelRepository) GetLabels(to interface{}) error {
	if result := r.getDb().Order("id").Find(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *LabelRepository) SaveLabel(label interface{}) error {
	result := r.getDb().Save(label)
	return result.Error
}

func (r *LabelRepository) GetLabelsByName(to interface{}, name string) error {
	if result := r.getDb().Where("name = ?", name).Find(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *LabelRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
