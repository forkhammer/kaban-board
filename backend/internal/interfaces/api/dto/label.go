package dto

import (
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type UpdateLabelRequest struct {
	AltName       *string `json:"altName"`
	BindingStatus *string `json:"bindingStatus"`
	Priority      *string `json:"priority"`
}

type LabelDto struct {
	Id            string  `json:"id"`
	Title         string  `json:"name"`
	Color         string  `json:"color"`
	TextColor     string  `json:"textColor"`
	AltName       *string `json:"altName"`
	BindingStatus *string `json:"bindingStatus"`
	Priority      *string `json:"priority"`
}

type KanbanLabelDto struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	AltName       *string `json:"altName"`
	BindingStatus *string `json:"bindingStatus"`
	Priority      *string `json:"priority"`
}

func SerializeLabel(label *domain.Label) *LabelDto {
	return &LabelDto{
		Id:            string(label.Id),
		Title:         label.Name,
		Color:         string(label.Color),
		TextColor:     string(label.TextColor),
		AltName:       label.AltName,
		BindingStatus: (*string)(label.BindingStatus),
		Priority:      (*string)(label.Priority),
	}
}

func SerializeLabels(labels []domain.Label) []LabelDto {
	result := utils.Map(labels, func(label domain.Label) LabelDto {
		return *SerializeLabel(&label)
	})
	return result
}

func SerializeKanbanLabel(label *usecases.KanbanLabel) *KanbanLabelDto {
	return &KanbanLabelDto{
		Id:            label.Name,
		Name:          label.Name,
		AltName:       label.AltName,
		BindingStatus: (*string)(label.BindingStatus),
		Priority:      (*string)(label.Priority),
	}
}

func SerializeKanbanLabels(labels []usecases.KanbanLabel) []KanbanLabelDto {
	result := utils.Map(labels, func(label usecases.KanbanLabel) KanbanLabelDto {
		return *SerializeKanbanLabel(&label)
	})
	return result
}
