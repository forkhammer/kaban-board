package usecases

import (
	"fmt"
	"log"
	"main/internal/app/interfaces"
	"main/internal/app/queries"
	app_services "main/internal/app/services"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"time"

	"github.com/getsentry/sentry-go"
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
	Planned     *bool
	Version     uint
}

type IssueBindingOrdering struct {
	Id    uint
	Order string
}

type IssueBindingOrderingSet = []IssueBindingOrdering

type CreateIssueBindingRequest struct {
	Title           string
	ProjectId       uint
	SprintId        uint
	AssigneeId      *uint
	CreateInTracker bool
	AccountId       domain.AccountId
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
	labelRepo         repo.LabelRepo                           `di.inject:"LabelRepository"`
	labelQuery        queries.LabelQuery                       `di.inject:"LabelQuery"`
	gitlab            interfaces.TaskTracker                   `di.inject:"gitlab"`
	gitlabAuthService interfaces.GitLabAuthServiceInterface    `di.inject:"GitLabAuthService"`
	sprintHub         interfaces.SprintHub                     `di.inject:"SprintHub"`
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
		u.commonQuery.OrderSpec("issue_bindings.\"order\", issue_bindings.assignee_id, issue_bindings.id"),
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

	if err := u.issueBindingRepo.Delete(domain.IssueBindingId(id)); err != nil {
		return 0, err
	}
	u.sprintHub.Broadcast(sprintId, interfaces.WSEvent{
		Type: interfaces.EventBindingDeleted,
		Data: interfaces.BindingDeletedEvent{ID: uint(id), SprintID: sprintId},
	})
	return sprintId, nil
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

	oldRelease := binding.Release

	var oldEpic *domain.Epic
	if binding.Epic != nil {
		oldEpic = binding.Epic
	}

	oldPriority := binding.Priority

	binding.EstimateDev = request.EstimateDev
	binding.EstimateQA = request.EstimateQA
	if request.BindStatus != nil {
		binding.BindStatus = *request.BindStatus
	}
	binding.Comment = request.Comment
	binding.Priority = request.Priority

	if request.Planned != nil {
		binding.Planned = request.Planned
	}

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

	releaseChanged := !app_services.EqualReleases(oldRelease, binding.Release)

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

	if releaseChanged && binding.Issue.IsExternal() {
		go uc.saveMilestoneToTracker(binding, account)
	}

	epicChanged := !app_services.EqualEpics(oldEpic, binding.Epic)

	if epicChanged && binding.Issue.IsExternal() {
		if oldEpic != nil {
			go uc.removeEpicLinkFromTracker(binding, oldEpic, account)
		}
		if binding.Epic != nil {
			go uc.saveEpicLinkToTracker(binding, account)
		}
	}

	priorityChanged := !app_services.EqualPriority(oldPriority, binding.Priority)
	if priorityChanged && binding.Issue.IsExternal() {
		go uc.savePriorityToTracker(binding, oldPriority, account)
	}

	_, err = uc.historyService.AddHistory(binding)
	if err != nil {
		return nil, err
	}

	uc.sprintHub.Broadcast(binding.Sprint.Id, interfaces.WSEvent{
		Type: interfaces.EventBindingUpdated,
		Data: interfaces.BindingUpdatedEvent{ID: binding.Id, AccountID: account.Id},
	})

	return binding, nil
}

func (uc *IssueBindingUseCases) saveMilestoneToTracker(binding *domain.IssueBinding, account *domain.Account) {
	token, err := uc.gitlabAuthService.GetValidToken(account.Id)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("saveMilestoneToTracker: failed to get gitlab token for account %d: %v", account.Id, err)
		return
	}
	var milestoneId *uint
	if binding.Release != nil {
		id := uint(binding.Release.Id)
		milestoneId = &id
	}
	if err := uc.gitlab.UpdateIssueMilestone(
		token.AccessToken,
		uint(binding.Issue.Project.Id),
		string(binding.Issue.Iid),
		milestoneId,
	); err != nil {
		sentry.CaptureException(err)
		log.Printf("saveMilestoneToTracker: failed to update milestone for issue %s: %v", binding.Issue.Iid, err)
	}
}

func (uc *IssueBindingUseCases) removeEpicLinkFromTracker(binding *domain.IssueBinding, epic *domain.Epic, account *domain.Account) {
	token, err := uc.gitlabAuthService.GetValidToken(account.Id)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("removeEpicLinkFromTracker: failed to get gitlab token for account %d: %v", account.Id, err)
		return
	}
	if err := uc.gitlab.RemoveIssueLink(
		token.AccessToken,
		uint(binding.Issue.Project.Id),
		string(binding.Issue.Iid),
		uint(epic.Project.Id),
		string(epic.Iid),
	); err != nil {
		sentry.CaptureException(err)
		log.Printf("removeEpicLinkFromTracker: failed to remove issue link for issue %s -> epic %s: %v", binding.Issue.Iid, epic.Iid, err)
	}
}

func (uc *IssueBindingUseCases) saveEpicLinkToTracker(binding *domain.IssueBinding, account *domain.Account) {
	token, err := uc.gitlabAuthService.GetValidToken(account.Id)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("saveEpicLinkToTracker: failed to get gitlab token for account %d: %v", account.Id, err)
		return
	}
	epic := binding.Epic
	if err := uc.gitlab.AddIssueLink(
		token.AccessToken,
		uint(binding.Issue.Project.Id),
		string(binding.Issue.Iid),
		uint(epic.Project.Id),
		string(epic.Iid),
	); err != nil {
		sentry.CaptureException(err)
		log.Printf("saveEpicLinkToTracker: failed to add issue link for issue %s -> epic %s: %v", binding.Issue.Iid, epic.Iid, err)
	}
}

func (uc *IssueBindingUseCases) savePriorityToTracker(binding *domain.IssueBinding, oldPriority *domain.IssueBindingPriority, account *domain.Account) {
	token, err := uc.gitlabAuthService.GetValidToken(account.Id)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("savePriorityToTracker: failed to get gitlab token for account %d: %v", account.Id, err)
		return
	}

	projectId := uint(binding.Issue.Project.Id)
	issueIid := string(binding.Issue.Iid)

	if oldPriority != nil {
		spec := uc.labelQuery.GetSpec(queries.LabelFilter{Priority: oldPriority})
		labels, err := uc.labelRepo.List(spec)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("savePriorityToTracker: failed to find old priority label: %v", err)
		} else if len(labels) > 0 {
			if err := uc.gitlab.RemoveIssueLabel(token.AccessToken, projectId, issueIid, labels[0].Name); err != nil {
				sentry.CaptureException(err)
				log.Printf("savePriorityToTracker: failed to remove label %s: %v", labels[0].Name, err)
			}
		}
	}

	if binding.Priority != nil {
		spec := uc.labelQuery.GetSpec(queries.LabelFilter{Priority: binding.Priority})
		labels, err := uc.labelRepo.List(spec)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("savePriorityToTracker: failed to find new priority label: %v", err)
			return
		}
		if len(labels) > 0 {
			if err := uc.gitlab.AddIssueLabel(token.AccessToken, projectId, issueIid, labels[0].Name); err != nil {
				sentry.CaptureException(err)
				log.Printf("savePriorityToTracker: failed to add label %s: %v", labels[0].Name, err)
			}
		}
	}
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
	if sprintId > 0 {
		uc.sprintHub.Broadcast(sprintId, interfaces.WSEvent{
			Type: interfaces.EventBindingOrdering,
			Data: ordering,
		})
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
		Planned:     original.Planned,
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

	uc.sprintHub.Broadcast(newBinding.Sprint.Id, interfaces.WSEvent{
		Type: interfaces.EventBindingCreated,
		Data: interfaces.BindingCreatedEvent{ID: newBinding.Id, AccountID: account.Id},
	})

	return newBinding, nil
}

func (uc *IssueBindingUseCases) MoveBinding(request MoveIssueBindingRequest, account *domain.Account) (*domain.IssueBinding, error) {
	if account == nil {
		return nil, domain_pkg.NewValidationError("нельзя переместить анонимным пользователем")
	}

	binding, err := uc.issueBindingRepo.Get(domain.IssueBindingId(request.Id))
	if err != nil {
		return nil, err
	}

	if !binding.CanManage(account) {
		return nil, domain_pkg.NewValidationError("нельзя переместить этим пользователем")
	}

	oldSprintId := binding.Sprint.Id

	sprint, err := uc.sprintRepo.Get(domain.SprintId(request.SprintId))
	if err != nil {
		return nil, err
	}

	if _, err = uc.historyService.AddDeleteHistory(binding, time.Now()); err != nil {
		return nil, err
	}

	binding.Sprint = sprint

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

	uc.sprintHub.Broadcast(oldSprintId, interfaces.WSEvent{
		Type: interfaces.EventBindingDeleted,
		Data: interfaces.BindingDeletedEvent{ID: uint(request.Id), SprintID: oldSprintId},
	})
	uc.sprintHub.Broadcast(binding.Sprint.Id, interfaces.WSEvent{
		Type: interfaces.EventBindingCreated,
		Data: interfaces.BindingCreatedEvent{ID: binding.Id, AccountID: account.Id},
	})

	return binding, nil
}

func (uc *IssueBindingUseCases) CreateIssueAndBinding(request CreateIssueBindingRequest) (*domain.IssueBinding, error) {
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

	var issue *domain.Issue
	if request.CreateInTracker {
		token, err := uc.gitlabAuthService.GetValidToken(request.AccountId)
		if err != nil {
			return nil, domain_pkg.NewValidationError("Нет авторизации в Gitlab", err)
		}
		issue, err = uc.gitlab.CreateIssue(token.AccessToken, request.Title, request.ProjectId, request.AssigneeId)
		if err != nil {
			return nil, err
		}
		issue, err = uc.issueRepo.Create(issue)
		if err != nil {
			return nil, err
		}
	} else {
		project, err := uc.projectRepo.Get(domain.ProjectId(request.ProjectId))
		if err != nil {
			return nil, err
		}

		externalId := domain.IssueExternalId(uuid.Must(uuid.NewV7()).String())
		issue = &domain.Issue{
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

	uc.sprintHub.Broadcast(binding.Sprint.Id, interfaces.WSEvent{
		Type: interfaces.EventBindingCreated,
		Data: interfaces.BindingCreatedEvent{ID: binding.Id, AccountID: request.AccountId},
	})

	return binding, nil
}
