package dto

import (
	domain "main/internal/domain/models"
)

type GetReleasesRequest struct {
	Search  *string `form:"search"`
	Project *uint   `form:"project"`
}

type ReleaseDto struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	WebPath string `json:"webPath"`
}

func SerializeRelease(release *domain.Release) *ReleaseDto {
	return &ReleaseDto{
		Id:      uint(release.Id),
		Title:   release.Title,
		WebPath: release.WebPath,
	}
}

func SerializeReleases(releases *[]domain.Release) *[]ReleaseDto {
	result := make([]ReleaseDto, len(*releases))
	for i, release := range *releases {
		result[i] = *SerializeRelease(&release)
	}
	return &result
}
