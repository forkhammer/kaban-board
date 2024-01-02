package rdb

import (
	"main/tools"

	"gorm.io/gorm"
)

type KVStoreRepository struct {
	connection tools.ConnectionInterface
}

func (r *KVStoreRepository) GetAll(to interface{}) error {
	if result := r.getDb().Find(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *KVStoreRepository) GetOrCreate(key string, to interface{}) error {
	if result := r.getDb().Where("key = ?", key).Attrs(to).FirstOrCreate(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *KVStoreRepository) Save(value interface{}) error {
	if result := r.getDb().Save(value); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *KVStoreRepository) Delete(key string) error {
	return r.getDb().Delete("key = ?", key).Error
}

func (r *KVStoreRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
