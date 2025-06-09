package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type IssueBindingRepository struct {
	conn       interfaces.ConnectionInterface `di.inject:"db"`
	sprintRepo *SprintRepository              `di.inject:"SprintRepository"`
	issueRepo  *IssueRepository               `di.inject:"IssueRepository"`
	userRepo   *UserRepository                `di.inject:"UserRepository"`
}

func (r *IssueBindingRepository) Get(id domain.IssueBindingId) (*domain.IssueBinding, error) {
	binding := &models.IssueBinding{}
	if err := r.getQuery().Where("id = ?", id).First(binding).Error; err != nil {
		return nil, err
	}
	model, err := r.toDomainIssueBinding(binding)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *IssueBindingRepository) List(spec repo.QuerySpec) (*[]domain.IssueBinding, error) {
	сolumns := make([]models.IssueBinding, 0)
	query := r.getQuery().Model(&models.IssueBinding{})

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

	return &domainIssueBindings, nil
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
	if err := r.conn.GetEngine().Save(model).Error; err != nil {
		return nil, err
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
		assignee = r.userRepo.toDomainUser(binding.Assignee)
	}

	domainBinding := &domain.IssueBinding{
		Id:          domain.IssueBindingId(binding.Id),
		Issue:       issue,
		Sprint:      sprint,
		EstimateDev: binding.EstimateDev,
		EstimateQA:  binding.EstimateQA,
		BindStatus:  domain.IssueBindingStatus(binding.BindStatus),
		Assignee:    assignee,
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
		AssigneeId:  assigneeId,
	}, nil
}

func (r *IssueBindingRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.IssueBinding{}).Preload("Sprint").Preload("Issue").Preload("Assignee")
}
