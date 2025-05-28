package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
)

type ColumnDto struct {
	Id     uint     `json:"id"`
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	TeamId *uint    `json:"team_id"`
	Team   *TeamDto `json:"team"`
	Order  *int     `json:"order"`
}

func SerializeColumn(column *domain.Column) *ColumnDto {
	var teamId *uint
	if column.Team != nil {
		val := uint(column.Team.Id)
		teamId = &val
	}
	var teamDto *TeamDto
	if column.Team != nil {
		teamDto = SerializeTeam(column.Team)
	}

	return &ColumnDto{
		Id:   uint(column.Id),
		Name: column.Name,
		Labels: utils.Map(column.Labels, func(labelId domain.LabelId) string {
			return string(labelId)
		}),
		TeamId: teamId,
		Team:   teamDto,
		Order:  column.Order,
	}
}

func SerializeColumns(columns *[]domain.Column) *[]ColumnDto {
	result := utils.Map(*columns, func(column domain.Column) ColumnDto {
		return *SerializeColumn(&column)
	})
	return &result
}
