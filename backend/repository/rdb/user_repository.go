package rdb

import (
	"main/tools"

	"gorm.io/gorm"
)

type UserRepository struct {
	connection tools.ConnectionInterface
}

func (r *UserRepository) GetUsers(to interface{}) error {
	result := r.getDb().Order("name").Find(to)
	return result.Error
}

func (r *UserRepository) GetVisibleUsers(to interface{}) error {
	result := r.getDb().Where("is_visible = true").Order("name").Find(to)
	return result.Error
}

func (r *UserRepository) GetOrCreate(to, query, attrs interface{}) error {
	if result := r.getDb().Where(query).Attrs(attrs).FirstOrCreate(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *UserRepository) GetUserBydId(to interface{}, id int) error {
	if result := r.getDb().Where("id = ?", id).First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *UserRepository) SaveUser(user interface{}) error {
	result := r.getDb().Save(user)
	return result.Error
}

func (r *UserRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
