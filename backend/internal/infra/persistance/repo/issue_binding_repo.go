package repo

import (
	"errors"
	"strconv"

	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type IssueBindingRepository struct {
	conn        interfaces.ConnectionInterface `di.inject:"db"`
	sprintRepo  *SprintRepository              `di.inject:"SprintRepository"`
	issueRepo   *IssueRepository               `di.inject:"IssueRepository"`
	userRepo    *UserRepository                `di.inject:"UserRepository"`
	releaseRepo *ReleaseRepository             `di.inject:"ReleaseRepository"`
	epicRepo    *EpicRepository                `di.inject:"EpicRepository"`
}

func (r *IssueBindingRepository) Get(id domain.IssueBindingId) (*domain.IssueBinding, error) {
	binding := &models.IssueBinding{}
	if err := r.getQuery().Where("issue_bindings.id = ?", id).First(binding).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_pkg.NewNotFoundError("IssueBinding", strconv.Itoa(int(id)), err)
		}
		return nil, err
	}
	model, err := r.toDomainIssueBinding(binding)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *IssueBindingRepository) List(spec repo.QuerySpec) ([]domain.IssueBinding, error) {
	сolumns := make([]models.IssueBinding, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&сolumns).Error; err != nil {
		return nil, err
	}

	domainIssueBindings := make([]domain.IssueBinding, len(сolumns))
	for i, IssueBinding := range сolumns {
		if elem, err := r.toDomainIssueBinding(&IssueBinding); err != nil {
			return nil, err
		} else {
			domainIssueBindings[i] = *elem
		}
	}

	return domainIssueBindings, nil
}

func (r *IssueBindingRepository) Count(spec repo.QuerySpec) (int, error) {
	var count int64
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return 0, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *IssueBindingRepository) Create(binding *domain.IssueBinding) (*domain.IssueBinding, error) {
	model, err := r.toIssueBinding(binding)
	if err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Create(model).Error; err != nil {
		return nil, err
	}

	return r.Get(domain.IssueBindingId(model.Id))
}

func (r *IssueBindingRepository) Update(binding *domain.IssueBinding) (*domain.IssueBinding, error) {
	model, err := r.toIssueBinding(binding)
	if err != nil {
		return nil, err
	}
	result := r.conn.GetEngine().Select("*").Omit("CreatedAt").Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	// Check for optimistic lock conflict (no rows updated)
	if result.RowsAffected == 0 {
		currentBinding, getErr := r.Get(binding.Id)
		if getErr != nil {
			return nil, getErr
		}
		return nil, domain_pkg.NewConflictError("IssueBinding", currentBinding, errors.New("resource was modified by another user"))
	}

	return r.Get(domain.IssueBindingId(model.Id))
}

func (r *IssueBindingRepository) Delete(id domain.IssueBindingId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.IssueBinding{}).Error
}

func (r *IssueBindingRepository) toDomainIssueBinding(binding *models.IssueBinding) (*domain.IssueBinding, error) {
	sprint, err := r.sprintRepo.toDomainSprint(binding.Sprint)
	if err != nil {
		return nil, err
	}

	issue, err := r.issueRepo.toDomainIssue(binding.Issue)
	if err != nil {
		return nil, err
	}

	var assignee *domain.User
	if binding.Assignee != nil {
		assignee = r.userRepo.toDomainUser(binding.Assignee, make(map[uint]*domain.Account))
	}

	var release *domain.Release
	if binding.Release != (*models.Release)(nil) {
		release, err = r.releaseRepo.toDomainRelease(binding.Release)
		if err != nil {
			return nil, err
		}
	}

	var epic *domain.Epic
	if binding.Epic != (*models.Epic)(nil) {
		epic, err = r.epicRepo.toDomainEpic(binding.Epic)
		if err != nil {
			return nil, err
		}
	}

		domainBinding := &domain.IssueBinding{
		Id:          domain.IssueBindingId(binding.Id),
		CreatedAt:   binding.CreatedAt,
		Issue:       issue,
		Sprint:      sprint,
		EstimateDev: binding.EstimateDev,
		EstimateQA:  binding.EstimateQA,
		BindStatus:  domain.IssueBindingStatus(binding.BindStatus),
		Priority:    (*domain.IssueBindingPriority)(binding.Priority),
		Assignee:    assignee,
		Comment:     binding.Comment,
		Release:     release,
		Epic:        epic,
		Order:       binding.Order,
		Planned:     binding.Planned,
		Version:     uint(binding.Version.Int64),
	}

	return domainBinding, domainBinding.Validate()
}

func (r *IssueBindingRepository) toIssueBinding(binding *domain.IssueBinding) (*models.IssueBinding, error) {
	var assigneeId *uint
	if binding.Assignee != nil {
		val := uint(binding.Assignee.Id)
		assigneeId = &val
	}

	return &models.IssueBinding{
		Id:          uint(binding.Id),
		IssueId:     uint(binding.Issue.Id),
		SprintId:    uint(binding.Sprint.Id),
		EstimateDev: binding.EstimateDev,
		EstimateQA:  binding.EstimateQA,
		BindStatus:  string(binding.BindStatus),
		Priority:    (*string)(binding.Priority),
		AssigneeId:  assigneeId,
		Comment:     binding.Comment,
		Order:       binding.Order,
		Planned:     binding.Planned,
		ReleaseId: func() *uint {
			if binding.Release != (*domain.Release)(nil) {
				val := uint(binding.Release.Id)
				return &val
			}
			return nil
		}(),
		EpicId: func() *uint {
			if binding.Epic != (*domain.Epic)(nil) {
				val := uint(binding.Epic.Id)
				return &val
			}
			return nil
		}(),
		Version: func() optimisticlock.Version {
			return optimisticlock.Version{Int64: int64(binding.Version), Valid: true}
		}(),
	}, nil
}

func (r *IssueBindingRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.IssueBinding{}).
		Joins("Sprint").
		Joins("Sprint.Team").
		Joins("Issue").
		Preload("Issue.Labels").
		Joins("Issue.Project").
		Joins("Issue.Project.Team").
		Joins("Issue.TaskType").
		Joins("Issue.Release").
		Joins("Assignee").
		Preload("Assignee.Groups").
		Joins("Release").
		Joins("Epic").
		Joins("Epic.Project")
}
