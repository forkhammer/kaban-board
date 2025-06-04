package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type QuarterDto struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func SerializeQuarter(quarter *domain.Quarter) *QuarterDto {
	return &QuarterDto{
		Id:    quarter.GetId(),
		Title: quarter.GetTitle(),
	}
}

func SerializeQuarters(quarters []domain.Quarter) []QuarterDto {
	return utils.Map(quarters, func(quarter domain.Quarter) QuarterDto {
		return *SerializeQuarter(&quarter)
	})
}
