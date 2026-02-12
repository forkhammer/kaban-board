package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

const DEFAULT_PAGE_SIZE int = 20

type SaveIssueRequest struct {
	Id          uint
	BindingId   *uint
	EstimateDev *uint
	EstimateQA  *uint
	BindStatus  *domain.IssueBindingStatus
	Assignee    *uint
	Comment     *string
	Priority    *domain.IssueBindingPriority
	ReleaseId   *uint
	EpicId      *uint
}

type IssuePage struct {
	Results []domain.Issue
	Count   int
	Page    int
	Limit   int
}

type IssueUseCases struct {
	issueQuery      queries.IssueQuery      `di.inject:"IssueQuery"`
	issueRepo       repo.IssueRepo          `di.inject:"IssueRepository"`
	sprintRepo      repo.SprintRepo         `di.inject:"SprintRepository"`
	userRepo        repo.UserRepo           `di.inject:"UserRepository"`
	releaseRepo     repo.ReleaseRepo        `di.inject:"ReleaseRepository"`
	epicRepo        repo.EpicRepo           `di.inject:"EpicRepository"`
	paginationQuery queries.PaginationQuery `di.inject:"PaginationQuery"`
}

func (u *IssueUseCases) GetIssues(filter *queries.IssueFilter, page int, limit int) (*IssuePage, error) {
	queryPage := page
	if page <= 0 {
		queryPage = 1
	}

	queryLimit := limit
	if queryLimit <= 0 {
		queryLimit = DEFAULT_PAGE_SIZE
	}
	query := u.paginationQuery.GetSpec(queryPage, queryLimit)

	if filter != nil {
		query = repo.And(query, u.issueQuery.GetSpec(*filter))
	}
	issues, err := u.issueRepo.List(query)
	if err != nil {
		return nil, err
	}
	count, err := u.issueRepo.Count(query)
	if err != nil {
		return nil, err
	}

	return &IssuePage{
		Results: issues,
		Count:   count,
		Page:    queryPage,
		Limit:   queryLimit,
	}, nil
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

	if err = issue.SetEstimateDev(request.EstimateDev); err != nil {
		return nil, err
	}
	if err = issue.SetEstimateQA(request.EstimateQA); err != nil {
		return nil, err
	}
	if request.BindStatus != nil {
		if err = issue.SetBindStatus(*request.BindStatus); err != nil {
			return nil, err
		}
	}
	if err = issue.SetComment(request.Comment); err != nil {
		return nil, err
	}
	if err = issue.SetPriority(request.Priority); err != nil {
		return nil, err
	}

	var release *domain.Release
	if request.ReleaseId != nil {
		if release, err = uc.releaseRepo.Get(domain.ReleaseId(*request.ReleaseId)); err != nil {
			return nil, err
		}
	}

	if err = issue.SetRelease(release); err != nil {
		return nil, err
	}

	var epic *domain.Epic
	if request.EpicId != nil {
		if epic, err = uc.epicRepo.Get(domain.EpicId(*request.EpicId)); err != nil {
			return nil, err
		}
	}

	if err = issue.SetEpic(epic); err != nil {
		return nil, err
	}

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
