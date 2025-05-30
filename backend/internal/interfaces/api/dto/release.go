package dto

import (
	domain "main/internal/domain/models"
)

type ReleaseDto struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	WebPath string `json:"webPath"`
}

func SerializeRelease(release *domain.Release) *ReleaseDto {
	return &ReleaseDto{
		Id:      string(release.Id),
		Title:   release.Title,
		WebPath: release.WebPath,
	}
}
