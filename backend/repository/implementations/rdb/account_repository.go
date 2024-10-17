package rdb

import (
	"main/repository"

	"gorm.io/gorm"
)

type AccountRepository struct {
	connection repository.ConnectionInterface
}

func (r *AccountRepository) CreateAccount(account interface{}) error {
	result := r.getDb().Create(account)

	return result.Error
}

func (r *AccountRepository) GetAccountByUsername(to interface{}, username string) error {
	if result := r.getDb().Where("username = ?", username).First(&to); result.Error == nil {
		return nil
	} else {
		return result.Error
	}
}

func (r *AccountRepository) GetAccountById(to interface{}, id uint) error {
	if result := r.getDb().Where("id = ?", id).First(&to); result.Error == nil {
		return nil
	} else {
		return result.Error
	}
}

func (r *AccountRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
