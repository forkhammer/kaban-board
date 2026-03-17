package usecases

import (
	"fmt"
	"main/internal/app/queries"
	app_services "main/internal/app/services"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"time"

	"github.com/google/uuid"
)

type IssueBindingPage struct {
	Results []domain.IssueBinding
	Count   int
	Page    int
	Limit   int
}

type SaveIssueBindingRequest struct {
	Id          uint
	Title       *string
	EstimateDev *uint
	EstimateQA  *uint
	BindStatus  *domain.IssueBindingStatus
	Assignee    *uint
	Comment     *string
	Priority    *domain.IssueBindingPriority
	ReleaseId   *uint
	EpicId      *uint
	Version     uint
}

type IssueBindingOrdering struct {
	Id    uint
	Order string
}

type IssueBindingOrderingSet = []IssueBindingOrdering

type CreateIssueBindingRequest struct {
	Title      string
	ProjectId  uint
	SprintId   uint
	AssigneeId *uint
}

type MoveIssueBindingRequest struct {
	Id       uint
	SprintId uint
}

type CopyIssueBindingRequest struct {
	Id       uint
	SprintId uint
}

type IssueBindingUseCases struct {
	commonQuery       queries.CommonQuery                      `di.inject:"CommonQuery"`
	issueBindingQuery queries.IssueBindingQuery                `di.inject:"IssueBindingQuery"`
	issueBindingRepo  repo.IssueBindingRepo                    `di.inject:"IssueBindingRepository"`
	issueRepo         repo.IssueRepo                           `di.inject:"IssueRepository"`
	projectRepo       repo.ProjectRepo                         `di.inject:"ProjectRepository"`
	sprintRepo        repo.SprintRepo                          `di.inject:"SprintRepository"`
	releaseRepo       repo.ReleaseRepo                         `di.inject:"ReleaseRepository"`
	epicRepo          repo.EpicRepo                            `di.inject:"EpicRepository"`
	userRepo          repo.UserRepo                            `di.inject:"UserRepository"`
	historyService    *app_services.IssueBindingHistoryService `di.inject:"IssueBindingHistoryService"`
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
		u.commonQuery.OrderSpec("issue_bindings.\"order\", issue_bindings.assignee_id, \"Issue\".created_at DESC"),
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

func (u *IssueBindingUseCases) DeleteBinding(id uint) (domain.SprintId, error) {
	binding, err := u.issueBindingRepo.Get(domain.IssueBindingId(id))
	if err != nil {
		return 0, err
	}

	sprintId := binding.Sprint.Id

	_, err = u.historyService.AddDeleteHistory(binding, time.Now())
	if err != nil {
		return 0, err
	}

	return sprintId, u.issueBindingRepo.Delete(domain.IssueBindingId(id))
}

func (uc *IssueBindingUseCases) SaveBinding(request SaveIssueBindingRequest, account *domain.Account) (*domain.IssueBinding, error) {
	if account == nil {
		return nil, fmt.Errorf("cannot be changed for an anonymous user")
	}

	binding, err := uc.issueBindingRepo.Get(domain.IssueBindingId(request.Id))
	if err != nil {
		return nil, err
	}

	if !binding.CanUpdate(account) {
		return nil, fmt.Errorf("cannot be changed for this user")
	}

	binding.EstimateDev = request.EstimateDev
	binding.EstimateQA = request.EstimateQA
	if request.BindStatus != nil {
		binding.BindStatus = *request.BindStatus
	}
	binding.Comment = request.Comment
	binding.Priority = request.Priority
	binding.Version = request.Version

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
		if binding.Assignee != nil && *request.Assignee != (uint)(binding.Assignee.Id) {
			if !binding.CanManage(account) {
				return nil, fmt.Errorf("cannot be managed for this user")
			}
		}
		assignee, err := uc.userRepo.Get(domain.UserId(*request.Assignee))
		if err != nil {
			return nil, err
		}
		binding.Assignee = assignee
	} else {
		if binding.Assignee != nil {
			if !binding.CanManage(account) {
				return nil, fmt.Errorf("cannot be managed for this user")
			}
		}
		binding.Assignee = nil
	}

	if request.Title != nil && binding.CanManage(account) && !binding.Issue.IsExternal() {
		issue := binding.Issue
		issue.Title = *request.Title

		if err := issue.Validate(); err != nil {
			return nil, err
		}

		_, err := uc.issueRepo.Update(issue)
		if err != nil {
			return nil, err
		}
	}

	if err = binding.Validate(); err != nil {
		return nil, err
	}

	binding, err = uc.issueBindingRepo.Update(binding)
	if err != nil {
		return nil, err
	}

	_, err = uc.historyService.AddHistory(binding)
	if err != nil {
		return nil, err
	}

	return binding, nil
}

func (uc *IssueBindingUseCases) SaveOrdering(ordering IssueBindingOrderingSet) (domain.SprintId, error) {
	var sprintId domain.SprintId
	for _, o := range ordering {
		binding, err := uc.issueBindingRepo.Get(domain.IssueBindingId(o.Id))
		if err != nil {
			return 0, err
		}
		if sprintId == 0 {
			sprintId = binding.Sprint.Id
		}
		binding.Order = o.Order
		if _, err := uc.issueBindingRepo.Update(binding); err != nil {
			return 0, err
		}
	}
	return sprintId, nil
}

func (uc *IssueBindingUseCases) CopyBinding(request CopyIssueBindingRequest, account *domain.Account) (*domain.IssueBinding, error) {
	if account == nil {
		return nil, domain_pkg.NewValidationError("нельзя скопировать анонимынм пользователем")
	}

	original, err := uc.issueBindingRepo.Get(domain.IssueBindingId(request.Id))
	if err != nil {
		return nil, err
	}

	if !original.CanManage(account) {
		return nil, domain_pkg.NewValidationError("нельяз скопировать этим пользователем")
	}

	sprint, err := uc.sprintRepo.Get(domain.SprintId(request.SprintId))
	if err != nil {
		return nil, err
	}

	newBinding := &domain.IssueBinding{
		Id:          0,
		Issue:       original.Issue,
		Sprint:      sprint,
		EstimateDev: original.EstimateDev,
		EstimateQA:  original.EstimateQA,
		BindStatus:  original.BindStatus,
		Priority:    original.Priority,
		Assignee:    original.Assignee,
		Comment:     original.Comment,
		Release:     original.Release,
		Epic:        original.Epic,
	}

	if err = newBinding.Validate(); err != nil {
		return nil, err
	}

	newBinding, err = uc.issueBindingRepo.Create(newBinding)
	if err != nil {
		return nil, err
	}

	_, err = uc.historyService.AddHistory(newBinding)
	if err != nil {
		return nil, err
	}

	return newBinding, nil
}

func (uc *IssueBindingUseCases) MoveBinding(request MoveIssueBindingRequest, account *domain.Account) (*domain.IssueBinding, domain.SprintId, error) {
	if account == nil {
		return nil, 0, domain_pkg.NewValidationError("нельзя переместить анонимным пользователем")
	}

	binding, err := uc.issueBindingRepo.Get(domain.IssueBindingId(request.Id))
	if err != nil {
		return nil, 0, err
	}

	if !binding.CanManage(account) {
		return nil, 0, domain_pkg.NewValidationError("нельзя переместить этим пользователем")
	}

	oldSprintId := binding.Sprint.Id

	sprint, err := uc.sprintRepo.Get(domain.SprintId(request.SprintId))
	if err != nil {
		return nil, 0, err
	}

	if _, err = uc.historyService.AddDeleteHistory(binding, time.Now()); err != nil {
		return nil, 0, err
	}

	binding.Sprint = sprint

	if err = binding.Validate(); err != nil {
		return nil, 0, err
	}

	binding, err = uc.issueBindingRepo.Update(binding)
	if err != nil {
		return nil, 0, err
	}

	_, err = uc.historyService.AddHistory(binding)
	if err != nil {
		return nil, 0, err
	}

	return binding, oldSprintId, nil
}

func (uc *IssueBindingUseCases) CreateIssueAndBinding(request CreateIssueBindingRequest) (*domain.IssueBinding, error) {
	project, err := uc.projectRepo.Get(domain.ProjectId(request.ProjectId))
	if err != nil {
		return nil, err
	}

	sprint, err := uc.sprintRepo.Get(domain.SprintId(request.SprintId))
	if err != nil {
		return nil, err
	}

	var assignee *domain.User
	if request.AssigneeId != nil {
		assignee, err = uc.userRepo.Get(domain.UserId(*request.AssigneeId))
		if err != nil {
			return nil, err
		}
	}

	externalId := domain.IssueExternalId(uuid.Must(uuid.NewV7()).String())

	issue := &domain.Issue{
		Id:         0,
		ExternalId: externalId,
		Iid:        "",
		Title:      request.Title,
		IssueType:  domain.IssueTypeIssue,
		Project:    *project,
		WebUrl:     "",
		Assignees:  []domain.User{},
		Labels:     []domain.Label{},
	}

	issue, err = uc.issueRepo.Create(issue)
	if err != nil {
		return nil, err
	}

	binding, err := issue.BindToSprint(sprint, assignee)
	if err != nil {
		return nil, err
	}

	binding, err = uc.issueBindingRepo.Create(binding)
	if err != nil {
		return nil, err
	}

	_, err = uc.historyService.AddHistory(binding)
	if err != nil {
		return nil, err
	}

	return binding, nil
}
