package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type SaveIssueRequest struct {
	Id          uint
	BindingId   *uint
	EstimateDev *uint
	EstimateQA  *uint
	BindStatus  *domain.IssueBindingStatus
	Assignee    *uint
	Comment     *string
}

type IssueUseCases struct {
	issueQuery queries.IssueQuery `di.inject:"IssueQuery"`
	issueRepo  repo.IssueRepo     `di.inject:"IssueRepository"`
	sprintRepo repo.SprintRepo    `di.inject:"SprintRepository"`
	userRepo   repo.UserRepo      `di.inject:"UserRepository"`
}

func (u *IssueUseCases) GetIssues(filter *queries.IssueFilter) (*[]domain.Issue, error) {
	var query repo.QuerySpec

	if filter != nil {
		query = u.issueQuery.GetSpec(*filter)
	}
	issues, err := u.issueRepo.List(query)
	if err != nil {
		return nil, err
	}

	return issues, nil
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

func (u *IssueUseCases) UnbindIssue(id uint, bindingId uint) (*domain.Issue, error) {
	issue, err := u.issueRepo.Get(domain.IssueId(id))
	if err != nil {
		return nil, err
	}

	err = issue.UnbindFromSprint(domain.IssueBindingId(bindingId))
	if err != nil {
		return nil, err
	}

	return u.issueRepo.Update(issue)
}

func (uc *IssueUseCases) SaveIssue(request SaveIssueRequest) (*domain.Issue, error) {
	issue, err := uc.issueRepo.Get(domain.IssueId(request.Id))
	if err != nil {
		return nil, err
	}

	if request.BindingId != nil {
		issue.SetContext((*domain.IssueBindingId)(request.BindingId))
	}

	issue.SetEstimateDev(request.EstimateDev)
	issue.SetEstimateQA(request.EstimateQA)
	if request.BindStatus != nil {
		issue.SetBindStatus(*request.BindStatus)
	}
	issue.SetComment(request.Comment)

	if request.Assignee != nil {
		assignee, err := uc.userRepo.Get((domain.UserId)(*request.Assignee))
		if err != nil {
			return nil, err
		}
		issue.SetAssignee(assignee)
	}

	issue, err = uc.issueRepo.Update(issue)
	if err != nil {
		return nil, err
	}

	return issue, nil

}
