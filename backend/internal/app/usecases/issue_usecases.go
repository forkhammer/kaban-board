package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type IssueUseCases struct {
	issueQuery queries.IssueQuery `di.inject:"IssueQuery"`
	issueRepo  repo.IssueRepo     `di.inject:"IssueRepository"`
}

func (u *IssueUseCases) GetIssues(filter *queries.IssueFilter) (*[]domain.Issue, error) {
	var query repo.QuerySpec
	if filter != nil {
		query = u.issueQuery.GetSpec(*filter)
	}
	return u.issueRepo.List(query)
}

func (u *IssueUseCases) GetIssue(id uint) (*domain.Issue, error) {
	return u.issueRepo.Get(domain.IssueId(id))
}
