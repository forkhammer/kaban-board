package rdb

import (
	"main/repository"

	"gorm.io/gorm"
)

type GroupRepository struct {
	connection repository.ConnectionInterface
}

func (r *GroupRepository) GetGroups(to interface{}) error {
	result := r.getDb().Find(to)
	return result.Error
}

func (r *GroupRepository) GetGroupById(to interface{}, id int) error {
	if result := r.getDb().Where("id = ?", id).First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *GroupRepository) GetGroupsByIds(to interface{}, ids []int) error {
	if result := r.getDb().Where("id IN ?", ids).Find(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *GroupRepository) SaveGroup(group interface{}) error {
	if result := r.getDb().Save(group); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *GroupRepository) CreateGroup(group interface{}) error {
	if result := r.getDb().Create(group); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *GroupRepository) DeleteGroup(group interface{}) error {
	result := r.getDb().Delete(group)
	return result.Error
}

func (r *GroupRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
