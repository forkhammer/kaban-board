package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type IssueBindingPage struct {
	Results []domain.IssueBinding
	Count   int
	Page    int
	Limit   int
}

type SaveIssueBindingRequest struct {
	Id          uint
	EstimateDev *uint
	EstimateQA  *uint
	BindStatus  *domain.IssueBindingStatus
	Assignee    *uint
	Comment     *string
	Priority    *domain.IssueBindingPriority
	ReleaseId   *uint
	EpicId      *uint
}

type IssueBindingUseCases struct {
	commonQuery       queries.CommonQuery       `di.inject:"CommonQuery"`
	issueBindingQuery queries.IssueBindingQuery `di.inject:"IssueBindingQuery"`
	issueBindingRepo  repo.IssueBindingRepo     `di.inject:"IssueBindingRepository"`
	releaseRepo       repo.ReleaseRepo          `di.inject:"ReleaseRepository"`
	epicRepo          repo.EpicRepo             `di.inject:"EpicRepository"`
	userRepo          repo.UserRepo             `di.inject:"UserRepository"`
}

func (u *IssueBindingUseCases) GetBindings(filter *queries.IssueBindingFilter, page int, limit int, account *domain.Account) (*IssueBindingPage, error) {
	queryPage := page
	if page <= 0 {
		queryPage = 1
	}

	queryLimit := limit
	if queryLimit <= 0 {
		queryLimit = DEFAULT_PAGE_SIZE
	}
	query := repo.And(
		u.commonQuery.PaginationSpec(queryPage, queryLimit),
		u.commonQuery.OrderSpec("Issue.created_at DESC"),
	)

	var filterQuery repo.QuerySpec
	if filter != nil {
		filterQuery = u.issueBindingQuery.GetSpec(*filter)
	}

	query = repo.And(query, filterQuery)
	issueBindings, err := u.issueBindingRepo.List(query)
	if err != nil {
		return nil, err
	}
	count, err := u.issueBindingRepo.Count(filterQuery)
	if err != nil {
		return nil, err
	}

	return &IssueBindingPage{
		Results: issueBindings,
		Count:   count,
		Page:    queryPage,
		Limit:   queryLimit,
	}, nil
}

func (u *IssueBindingUseCases) GetBinding(id uint, account *domain.Account) (*domain.IssueBinding, error) {
	return u.issueBindingRepo.Get(domain.IssueBindingId(id))
}

func (u *IssueBindingUseCases) DeleteBinding(id uint) error {
	return u.issueBindingRepo.Delete(domain.IssueBindingId(id))
}

func (uc *IssueBindingUseCases) SaveBinding(request SaveIssueBindingRequest) (*domain.IssueBinding, error) {
	binding, err := uc.issueBindingRepo.Get(domain.IssueBindingId(request.Id))
	if err != nil {
		return nil, err
	}

	binding.EstimateDev = request.EstimateDev
	binding.EstimateQA = request.EstimateQA
	if request.BindStatus != nil {
		binding.BindStatus = *request.BindStatus
	}
	binding.Comment = request.Comment
	binding.Priority = request.Priority

	if request.ReleaseId != nil {
		release, err := uc.releaseRepo.Get(domain.ReleaseId(*request.ReleaseId))
		if err != nil {
			return nil, err
		}
		binding.Release = release
	} else {
		binding.Release = nil
	}

	if request.EpicId != nil {
		epic, err := uc.epicRepo.Get(domain.EpicId(*request.EpicId))
		if err != nil {
			return nil, err
		}
		binding.Epic = epic
	} else {
		binding.Epic = nil
	}

	if request.Assignee != nil {
		assignee, err := uc.userRepo.Get(domain.UserId(*request.Assignee))
		if err != nil {
			return nil, err
		}
		binding.Assignee = assignee
	} else {
		binding.Assignee = nil
	}

	if err = binding.Validate(); err != nil {
		return nil, err
	}

	binding, err = uc.issueBindingRepo.Update(binding)
	if err != nil {
		return nil, err
	}

	return binding, nil
}
