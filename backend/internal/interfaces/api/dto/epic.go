package dto

import (
	domain "main/internal/domain/models"
)

type GetEpicsRequest struct {
	Search  *string `form:"search"`
	Project *uint   `form:"project"`
}

type EpicDto struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

func SerializeEpic(release *domain.Epic) *EpicDto {
	return &EpicDto{
		Id:    uint(release.Id),
		Title: release.Title,
	}
}

func SerializeEpics(releases []domain.Epic) []EpicDto {
	result := make([]EpicDto, len(releases))
	for i, release := range releases {
		result[i] = *SerializeEpic(&release)
	}
	return result
}
