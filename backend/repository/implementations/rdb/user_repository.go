package rdb

import (
	"main/repository"
	"main/repository/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	connection repository.ConnectionInterface
}

func (r *UserRepository) GetUsers(to *[]models.User) error {
	result := r.getDb().Preload("Groups").Order("name").Find(to)
	return result.Error
}

func (r *UserRepository) GetVisibleUsers(to *[]models.User) error {
	result := r.getDb().Preload("Groups").Where("is_visible = true").Order("name").Find(to)
	return result.Error
}

func (r *UserRepository) GetOrCreate(to *models.User, query, attrs interface{}) error {
	if result := r.getDb().Where(query).Attrs(attrs).FirstOrCreate(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *UserRepository) GetUserBydId(to *models.User, id int) error {
	if result := r.getDb().Where("id = ?", id).Preload("Groups").First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *UserRepository) SaveUser(user *models.User) error {
	if result := r.getDb().Save(user); result.Error != nil {
		return result.Error
	}

	if err := r.getDb().Model(user).Association("Groups").Replace(user.Groups); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
