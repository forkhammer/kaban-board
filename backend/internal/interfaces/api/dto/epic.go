package dto

import (
	domain "main/internal/domain/models"
)

type GetEpicsRequest struct {
	Search  *string `form:"search"`
	Project *uint   `form:"project"`
}

type EpicDto struct {
	Id      uint   `json:"id"`
	Iid     string `json:"iid"`
	Title   string `json:"title"`
	Project string `json:"project"`
}

func SerializeEpic(release *domain.Epic) *EpicDto {
	return &EpicDto{
		Id:      uint(release.Id),
		Iid:     string(release.Iid),
		Title:   release.Title,
		Project: release.Project.Name,
	}
}

func SerializeEpics(releases []domain.Epic) []EpicDto {
	result := make([]EpicDto, len(releases))
	for i, release := range releases {
		result[i] = *SerializeEpic(&release)
	}
	return result
}
