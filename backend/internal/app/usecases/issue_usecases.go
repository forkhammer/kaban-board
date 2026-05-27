package usecases

import (
	"main/internal/app/queries"
	app_services "main/internal/app/services"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

const DEFAULT_PAGE_SIZE int = 20

type IssuePage struct {
	Results []domain.Issue
	Count   int
	Page    int
	Limit   int
}

type BindIssuesResult struct {
	Bindings []domain.IssueBinding
	Errors   []string
}

type IssueUseCases struct {
	issueQuery        queries.IssueQuery                       `di.inject:"IssueQuery"`
	issueRepo         repo.IssueRepo                           `di.inject:"IssueRepository"`
	issueBindingRepo  repo.IssueBindingRepo                    `di.inject:"IssueBindingRepository"`
	issueBindingQuery queries.IssueBindingQuery                `di.inject:"IssueBindingQuery"`
	sprintRepo        repo.SprintRepo                          `di.inject:"SprintRepository"`
	commonQuery       queries.CommonQuery                      `di.inject:"CommonQuery"`
	userRepo          repo.UserRepo                            `di.inject:"UserRepository"`
	historyService    *app_services.IssueBindingHistoryService `di.inject:"IssueBindingHistoryService"`
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
	query := repo.And(
		u.commonQuery.PaginationSpec(queryPage, queryLimit),
		u.commonQuery.OrderSpec("issues.created_at DESC"),
	)

	var filterQuery repo.QuerySpec
	if filter != nil {
		filterQuery = u.issueQuery.GetSpec(*filter)
	}

	query = repo.And(query, filterQuery)
	issues, err := u.issueRepo.List(query)
	if err != nil {
		return nil, err
	}
	count, err := u.issueRepo.Count(filterQuery)
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

func (u *IssueUseCases) BindIssue(id uint, sprintId uint, assigneeId uint) (*domain.IssueBinding, error) {
	issue, err := u.issueRepo.Get(domain.IssueId(id))
	if err != nil {
		return nil, err
	}

	sprint, err := u.sprintRepo.Get(domain.SprintId(sprintId))
	if err != nil {
		return nil, err
	}

	var assignee *domain.User
	if assigneeId > 0 {
		assignee, err = u.userRepo.Get(domain.UserId(assigneeId))
		if err != nil {
			return nil, err
		}

		count, err := u.issueBindingRepo.Count(u.issueBindingQuery.GetSpec(queries.IssueBindingFilter{
			Issues:     []uint{id},
			SprintId:   &sprintId,
			AssigneeId: &assigneeId,
		}))
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, domain_pkg.NewValidationError("Задача уже привязана к спринту с выбранным исполнителем")
		}

		if len(assignee.Groups) > 0 {
			groupIds := make([]uint, len(assignee.Groups))
			for i, g := range assignee.Groups {
				groupIds[i] = uint(g.Id)
			}
			count, err = u.issueBindingRepo.Count(u.issueBindingQuery.GetSpec(queries.IssueBindingFilter{
				Issues:           []uint{id},
				SprintId:         &sprintId,
				AssigneeGroupIds: groupIds,
			}))
			if err != nil {
				return nil, err
			}
			if count > 0 {
				return nil, domain_pkg.NewValidationError("Задача уже привязана к спринту с пользователем из той же группы")
			}
		}
	}

	binding, err := issue.BindToSprint(sprint, assignee)
	if err != nil {
		return nil, err
	}
	binding, err = u.issueBindingRepo.Create(binding)
	if err != nil {
		return nil, err
	}

	_, err = u.historyService.AddHistory(binding)
	if err != nil {
		return nil, err
	}

	return binding, nil
}

func (u *IssueUseCases) BindIssues(ids []uint, sprintId uint, assigneeId uint) (*BindIssuesResult, error) {
	result := BindIssuesResult{
		Bindings: make([]domain.IssueBinding, 0),
		Errors:   make([]string, 0),
	}
	for _, id := range ids {
		binding, err := u.BindIssue(id, sprintId, assigneeId)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Bindings = append(result.Bindings, *binding)
		}
	}
	return &result, nil
}
