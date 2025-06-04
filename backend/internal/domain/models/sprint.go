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
	IsCompleted  bool
	Id           SprintId
	HoursPerUser uint `validate:"required,gt=0"`
	Title        string
	StartDate    time.Time `validate:"required"`
	EndDate      time.Time
	Team         Team `validate:"required"`
}

func NewSprint(id SprintId, title string, startDate time.Time, endDate time.Time, team Team, isCompleted bool, hoursPerUser uint) (*Sprint, error) {
	sprint := &Sprint{
		Id:           id,
		Title:        title,
		StartDate:    startDate,
		EndDate:      endDate,
		Team:         team,
		IsCompleted:  isCompleted,
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

func (s *Sprint) GetStatus() SprintStatus {
	if s.IsCompleted {
		return SprintStatusCompleted
	}
	if s.StartDate.Before(time.Now()) {
		return SprintStatusRunning
	}
	return SprintStatusWaiting
}

func (s *Sprint) SetCompleted() error {
	if s.GetStatus() != SprintStatusRunning {
		return fmt.Errorf("Нельзя завершить не запущенный спринт")
	}

	s.IsCompleted = true
	return nil
}

func (s *Sprint) GetQuarter() *Quarter {
	year := s.StartDate.Year()
	month := s.StartDate.Month()
	val, _ := NewQuarter(year, int((month-1)/3+1))
	return val
}
