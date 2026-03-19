package interfaces

import (
	domain "main/internal/domain/models"
)

type TaskTracker interface {
	GetProjects() ([]domain.Project, error)
	GetUsers() ([]domain.User, error)
	GetIssues() ([]domain.Issue, error)
	GetLabels() ([]domain.Label, error)
	GetReleases() ([]domain.Release, error)
	CreateIssue(userToken string, title string, projectId uint, assigneeId *uint) (*domain.Issue, error)
}
