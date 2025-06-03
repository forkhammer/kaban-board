package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
	"time"
)

type SprintDto struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Team         TeamDto `json:"team"`
	Status       string  `json:"status"`
	HoursPerUser int     `json:"hours_per_user"`
}

type CreateSprintRequest struct {
	Title        string    `json:"title"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HoursPerUser uint      `json:"hours_per_user"`
	TeamId       uint      `json:"team_id"`
}

type UpdateSprintRequest struct {
	Title        string    `json:"title"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HoursPerUser uint      `json:"hours_per_user"`
	TeamId       uint      `json:"team_id"`
}

func SerializeSprint(sprint *domain.Sprint) *SprintDto {
	return &SprintDto{
		Id:           int(sprint.Id),
		Title:        sprint.GetTitle(),
		StartDate:    sprint.StartDate.String(),
		EndDate:      sprint.EndDate.String(),
		Team:         *SerializeTeam(&sprint.Team),
		Status:       string(sprint.GetStatus()),
		HoursPerUser: int(sprint.HoursPerUser),
	}
}

func SerializeSprints(sprints []domain.Sprint) []SprintDto {
	return utils.Map(sprints, func(sprint domain.Sprint) SprintDto {
		return *SerializeSprint(&sprint)
	})
}
