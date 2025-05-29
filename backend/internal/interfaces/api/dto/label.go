package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type UpdateLabelRequest struct {
	AltName *string `json:"altName"`
}

type LabelDto struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	AltName *string `json:"altName"`
}

func SerializeLabel(label *domain.Label) *LabelDto {
	return &LabelDto{
		Id:      string(label.Id),
		Title:   label.Name,
		AltName: label.AltName,
	}
}

func SerializeLabels(labels *[]domain.Label) *[]LabelDto {
	result := utils.Map(*labels, func(label domain.Label) LabelDto {
		return *SerializeLabel(&label)
	})
	return &result
}
