package models

type AccountId uint

type Account struct {
	IsActive bool
	Id       AccountId
	Username string
	Password string
	Name     string
}
