package dto

import (
	"main/internal/app/label_usecases"
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type UpdateLabelRequest struct {
	AltName *string `json:"altName"`
}

type LabelDto struct {
	Id        string  `json:"id"`
	Title     string  `json:"name"`
	Color     string  `json:"color"`
	TextColor string  `json:"textColor"`
	AltName   *string `json:"altName"`
}

type KanbanLabelDto struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	AltName *string `json:"altName"`
}

func SerializeLabel(label *domain.Label) *LabelDto {
	return &LabelDto{
		Id:        string(label.Id),
		Title:     label.Name,
		Color:     string(label.Color),
		TextColor: string(label.TextColor),
		AltName:   label.AltName,
	}
}

func SerializeLabels(labels *[]domain.Label) *[]LabelDto {
	result := utils.Map(*labels, func(label domain.Label) LabelDto {
		return *SerializeLabel(&label)
	})
	return &result
}

func SerializeKanbanLabel(label *label_usecases.KanbanLabel) *KanbanLabelDto {
	return &KanbanLabelDto{
		Id:      label.Name,
		Name:    label.Name,
		AltName: label.AltName,
	}
}

func SerializeKanbanLabels(labels *[]label_usecases.KanbanLabel) *[]KanbanLabelDto {
	result := utils.Map(*labels, func(label label_usecases.KanbanLabel) KanbanLabelDto {
		return *SerializeKanbanLabel(&label)
	})
	return &result
}
