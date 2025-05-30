package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type GroupDto struct {
	Id   uint   `json:"id"`
	Name string `json:"title"`
}

func SerializeGroups(groups []domain.Group) []GroupDto {
	result := utils.Map(groups, func(group domain.Group) GroupDto {
		return *SerializeGroup(&group)
	})
	return result
}

func SerializeGroup(group *domain.Group) *GroupDto {
	return &GroupDto{
		Id:   uint(group.Id),
		Name: group.Name,
	}
}
