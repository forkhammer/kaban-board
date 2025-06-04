package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
	"time"
)

type SprintDto struct {
	Id           int        `json:"id"`
	Title        string     `json:"title"`
	StartDate    string     `json:"start_date"`
	EndDate      string     `json:"end_date"`
	Team         TeamDto    `json:"team"`
	Status       string     `json:"status"`
	HoursPerUser int        `json:"hours_per_user"`
	Quarter      QuarterDto `json:"quarter"`
}

type GetSprintsRequest struct {
	Team    *int    `form:"team"`
	Quarter *string `form:"quarter"`
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
		StartDate:    sprint.StartDate.Format("2006-01-02"),
		EndDate:      sprint.EndDate.Format("2006-01-02"),
		Team:         *SerializeTeam(&sprint.Team),
		Status:       string(sprint.GetStatus()),
		HoursPerUser: int(sprint.HoursPerUser),
		Quarter:      *SerializeQuarter(sprint.GetQuarter()),
	}
}

func SerializeSprints(sprints []domain.Sprint) []SprintDto {
	return utils.Map(sprints, func(sprint domain.Sprint) SprintDto {
		return *SerializeSprint(&sprint)
	})
}
