package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type SprintStatus string

const (
	SprintStatusRunning   SprintStatus = "running"
	SprintStatusCompleted SprintStatus = "completed"
	SprintStatusWaiting   SprintStatus = "waiting"
)

type SprintId int

type Sprint struct {
	Id           SprintId
	HoursPerUser uint `validate:"required,gt=0"`
	Title        string
	StartDate    time.Time `validate:"required"`
	EndDate      time.Time
	Team         Team `validate:"required"`
	Status       SprintStatus
}

func NewSprint(id SprintId, title string, startDate time.Time, endDate time.Time, team Team, status SprintStatus, hoursPerUser uint) (*Sprint, error) {
	sprint := &Sprint{
		Id:           id,
		Title:        title,
		StartDate:    startDate,
		EndDate:      endDate,
		Team:         team,
		Status:       status,
		HoursPerUser: hoursPerUser,
	}
	return sprint, sprint.Validate()
}

func (s *Sprint) Validate() error {
	validator := validator.New()
	if err := validator.Struct(s); err != nil {
		return err
	}

	if s.StartDate.After(s.EndDate) {
		return fmt.Errorf("Дата начала должна быть раньше даты окончания")
	}

	return nil
}

func (s *Sprint) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	return fmt.Sprintf("Спринт %s - %s", s.StartDate.Format("02.01.2006"), s.EndDate.Format("02.01.2006"))
}

func (s *Sprint) Run() error {
	if s.Status == SprintStatusRunning {
		return fmt.Errorf("Нельзя запустить этот спринт")
	}

	if s.StartDate.After(time.Now()) {
		return fmt.Errorf("Нельзя запустить спринт в будущем")
	}

	s.Status = SprintStatusRunning
	return nil
}

func (s *Sprint) Complete() error {
	if s.Status != SprintStatusRunning {
		return fmt.Errorf("Нельзя завершить не запущенный спринт")
	}

	s.Status = SprintStatusCompleted
	return nil
}

func (s *Sprint) GetQuarter() *Quarter {
	return NewQuarterFromDate(s.StartDate)
}

func (s *Sprint) CanDelete() bool {
	return s.Status == SprintStatusWaiting
}

func (s *Sprint) CanRun() bool {
	return (s.Status == SprintStatusWaiting || s.Status == SprintStatusCompleted) && s.StartDate.Before(time.Now()) && s.EndDate.After(time.Now())
}

func (s *Sprint) CanComplete() bool {
	return s.Status == SprintStatusRunning
}
