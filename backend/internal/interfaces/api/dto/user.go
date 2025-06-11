package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type SetUserVisibilityRequest struct {
	Visible bool `json:"visible"`
}

type SetUserGroupsRequest struct {
	Groups []uint `json:"groups"`
}

type GetUsersRequest struct {
	Search *string `form:"search"`
	TeamId *uint   `form:"team_id"`
}

type UserDto struct {
	Id        uint       `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	AvatarUrl string     `json:"avatar_url"`
	IsVisible bool       `json:"is_visible"`
	Groups    []GroupDto `json:"groups"`
}

func SerializeUsers(users []domain.User) []UserDto {
	result := utils.Map(users, func(user domain.User) UserDto {
		return *SerializeUser(&user)
	})
	return result
}

func SerializeUser(user *domain.User) *UserDto {
	return &UserDto{
		Id:        uint(user.Id),
		Name:      user.Name,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		IsVisible: user.IsVisible,
		Groups:    SerializeGroups(user.Groups),
	}
}
