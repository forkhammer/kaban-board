package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type IssueBindingHistoryRepository struct {
	conn     interfaces.ConnectionInterface `di.inject:"db"`
	userRepo *UserRepository                `di.inject:"UserRepository"`
}

func (r *IssueBindingHistoryRepository) Get(id domain.IssueBindingHistoryId) (*domain.IssueBindingHistory, error) {
	history := &models.IssueBindingHistory{}
	if err := r.getQuery().Where("issue_binding_histories.id = ?", id).First(history).Error; err != nil {
		return nil, err
	}
	model, err := r.toDomainIssueBindingHistory(history)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *IssueBindingHistoryRepository) List(spec repo.QuerySpec) ([]domain.IssueBindingHistory, error) {
	histories := make([]models.IssueBindingHistory, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&histories).Error; err != nil {
		return nil, err
	}

	domainHistories := make([]domain.IssueBindingHistory, len(histories))
	for i, history := range histories {
		if elem, err := r.toDomainIssueBindingHistory(&history); err != nil {
			return nil, err
		} else {
			domainHistories[i] = *elem
		}
	}

	return domainHistories, nil
}

func (r *IssueBindingHistoryRepository) Count(spec repo.QuerySpec) (int, error) {
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

func (r *IssueBindingHistoryRepository) Create(history *domain.IssueBindingHistory) (*domain.IssueBindingHistory, error) {
	model, err := r.toIssueBindingHistory(history)
	if err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Create(model).Error; err != nil {
		return nil, err
	}

	return r.Get(domain.IssueBindingHistoryId(model.ID))
}

func (r *IssueBindingHistoryRepository) Delete(id domain.IssueBindingHistoryId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.IssueBindingHistory{}).Error
}

func (r *IssueBindingHistoryRepository) toDomainIssueBindingHistory(history *models.IssueBindingHistory) (*domain.IssueBindingHistory, error) {
	var assignee *domain.User
	if history.Assignee != nil {
		assignee = r.userRepo.toDomainUser(history.Assignee, make(map[uint]*domain.Account))
	}

	domainHistory := &domain.IssueBindingHistory{
		Id:             domain.IssueBindingHistoryId(history.ID),
		IssueBindingId: domain.IssueBindingId(history.IssueBindingId),
		IssueId:        domain.IssueId(history.IssueId),
		EstimateDev:    history.EstimateDev,
		EstimateQA:     history.EstimateQA,
		BindStatus:     domain.IssueBindingStatus(history.BindStatus),
		Priority:       (*domain.IssueBindingPriority)(history.Priority),
		Assignee:       assignee,
		RemovedAt:      history.RemovedAt,
		CreatedAt:      history.CreatedAt,
	}

	return domainHistory, domainHistory.Validate()
}

func (r *IssueBindingHistoryRepository) toIssueBindingHistory(history *domain.IssueBindingHistory) (*models.IssueBindingHistory, error) {
	var assigneeId *uint
	if history.Assignee != nil {
		val := uint(history.Assignee.Id)
		assigneeId = &val
	}

	return &models.IssueBindingHistory{
		IssueBindingId: uint(history.IssueBindingId),
		IssueId:        uint(history.IssueId),
		EstimateDev:    history.EstimateDev,
		EstimateQA:     history.EstimateQA,
		BindStatus:     string(history.BindStatus),
		Priority:       (*string)(history.Priority),
		AssigneeId:     assigneeId,
		RemovedAt:      history.RemovedAt,
	}, nil
}

func (r *IssueBindingHistoryRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.IssueBindingHistory{}).
		Joins("Issue").
		Joins("IssueBinding").
		Joins("Assignee").
		Preload("Assignee.Groups")
}
