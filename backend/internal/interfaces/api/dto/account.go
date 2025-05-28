package dto

import (
	domain "main/internal/domain/models"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AccountDto struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type ActiveUserResponse struct {
	User *AccountDto `json:"user"`
}

func NewActiveUserResponse(account *domain.Account) ActiveUserResponse {
	var accountDto *AccountDto

	if account != nil {
		accountDto = &AccountDto{
			Id:       uint(account.Id),
			Username: account.Username,
			Name:     account.Name,
			IsActive: account.IsActive,
		}
	}

	return ActiveUserResponse{
		User: accountDto,
	}
}
