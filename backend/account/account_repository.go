package account

import "main/db"

type AccountRepository struct{}

func (r *AccountRepository) CreateAccount(account *Account) (*Account, error) {
	result := db.DefaultConnection.Db.Create(account)

	return account, result.Error
}

func (r *AccountRepository) GetAccountByUsername(username string) (*Account, error) {
	var account Account

	if result := db.DefaultConnection.Db.Where(Account{Username: username}).First(&account); result.Error == nil {
		return &account, nil
	} else {
		return nil, result.Error
	}
}

func (r *AccountRepository) GetAccountById(id uint) (*Account, error) {
	var account Account

	if result := db.DefaultConnection.Db.Where(Account{Id: id}).First(&account); result.Error == nil {
		return &account, nil
	} else {
		return nil, result.Error
	}
}
