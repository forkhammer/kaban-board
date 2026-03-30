package app_services

import domain "main/internal/domain/models"

func EqualReleases(a, b *domain.Release) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Id == b.Id
}

func EqualEpics(a, b *domain.Epic) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Id == b.Id
}

func EqualPriority(a, b *domain.IssueBindingPriority) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
