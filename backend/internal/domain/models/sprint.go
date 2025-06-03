package models

import (
	"fmt"
	"time"
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
	HoursPerUser uint
	Title        string
	StartDate    time.Time
	EndDate      time.Time
	Team         Team
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
	if s.StartDate.After(s.EndDate) {
		return fmt.Errorf("Дата начала должна быть раньше даты окончания")
	}

	return nil
}

func (s *Sprint) GetTitle() string {
	if s.Title != "" {
		return s.Title
	}

	return fmt.Sprintf("Спринт %s - %s", s.StartDate.String(), s.EndDate.String())
}

func (s *Sprint) GetStatus() SprintStatus {
	if s.IsCompleted {
		return SprintStatusCompleted
	}
	if s.StartDate.After(time.Now()) {
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
