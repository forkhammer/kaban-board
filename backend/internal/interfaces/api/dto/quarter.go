package dto

import (
	domain "main/internal/domain/models"
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
