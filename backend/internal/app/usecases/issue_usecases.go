package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type IssueUseCases struct {
	issueQuery queries.IssueQuery `di.inject:"IssueQuery"`
	issueRepo  repo.IssueRepo     `di.inject:"IssueRepository"`
	sprintRepo repo.SprintRepo    `di.inject:"SprintRepository"`
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

func (u *IssueUseCases) BindIssue(id uint, sprintId uint) (*domain.Issue, error) {
	issue, err := u.issueRepo.Get(domain.IssueId(id))
	if err != nil {
		return nil, err
	}

	sprint, err := u.sprintRepo.Get(domain.SprintId(sprintId))
	if err != nil {
		return nil, err
	}

	err = issue.BindToSprint(sprint)
	if err != nil {
		return nil, err
	}
	return u.issueRepo.Update(issue)
}
