package dto

import (
	domain "main/internal/domain/models"
	"main/pkg/utils"
	"time"
)

type SprintDto struct {
	Id           int        `json:"id"`
	Title        string     `json:"title"`
	Represent    string     `json:"represent"`
	StartDate    string     `json:"start_date"`
	EndDate      string     `json:"end_date"`
	TeamId       int        `json:"team_id"`
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
	Title        string `json:"title"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	HoursPerUser uint   `json:"hours_per_user"`
	TeamId       uint   `json:"team_id"`
}

func (r *CreateSprintRequest) Validate() error {
	_, err := r.GetStartDate()
	if err != nil {
		return err
	}
	_, err = r.GetEndDate()
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateSprintRequest) GetStartDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.StartDate)
}

func (r *CreateSprintRequest) GetEndDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.EndDate)
}

type UpdateSprintRequest struct {
	CreateSprintRequest
}

func SerializeSprint(sprint *domain.Sprint) *SprintDto {
	return &SprintDto{
		Id:           int(sprint.Id),
		Title:        sprint.Title,
		Represent:    sprint.GetTitle(),
		StartDate:    sprint.StartDate.Format("2006-01-02"),
		EndDate:      sprint.EndDate.Format("2006-01-02"),
		TeamId:       int(sprint.Team.Id),
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
