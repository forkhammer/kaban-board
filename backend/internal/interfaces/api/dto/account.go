package dto

import (
	domain "main/internal/domain/models"
)

type AccountDto struct {
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	IsActive     bool   `json:"isActive"`
	AvatarURL    string `json:"avatarUrl,omitempty"`
	AuthProvider string `json:"authProvider,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse = AccountDto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ActiveUserResponse struct {
	User *AccountDto `json:"user"`
}

func NewActiveUserResponse(account *domain.Account) ActiveUserResponse {
	var accountDto *AccountDto

	if account != nil {
		accountDto = SerializeAccount(account)
	}

	return ActiveUserResponse{
		User: accountDto,
	}
}

func SerializeAccount(account *domain.Account) *AccountDto {
	return &AccountDto{
		Id:           uint(account.Id),
		Username:     account.Username,
		Name:         account.Name,
		IsActive:     account.IsActive,
		AvatarURL:    account.AvatarURL,
		AuthProvider: string(account.AuthProvider),
	}
}
