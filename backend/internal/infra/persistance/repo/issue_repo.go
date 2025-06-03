package repo

import (
	"errors"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IssueRepository struct {
	conn        interfaces.ConnectionInterface `di.inject:"db"`
	userRepo    *UserRepository                `di.inject:"UserRepository"`
	projectRepo *ProjectRepository             `di.inject:"ProjectRepository"`
	releaseRepo *ReleaseRepository             `di.inject:"ReleaseRepository"`
	labelRepo   *LabelRepository               `di.inject:"LabelRepository"`
}

func (r *IssueRepository) Get(id domain.IssueId) (*domain.Issue, error) {
	issue := &models.Issue{}
	if err := r.getQuery().Where("id = ?", id).First(issue).Error; err != nil {
		return nil, err
	}
	return r.toDomainIssue(issue)
}

func (r *IssueRepository) List(spec repo.QuerySpec) (*[]domain.Issue, error) {
	issues := make([]models.Issue, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&issues).Error; err != nil {
		return nil, err
	}

	domainIssues := make([]domain.Issue, len(issues))
	for i, issue := range issues {
		model, err := r.toDomainIssue(&issue)
		if err != nil {
			return nil, err
		}
		domainIssues[i] = *model
	}

	return &domainIssues, nil
}

func (r *IssueRepository) Create(issue *domain.Issue) (*domain.Issue, error) {
	model := r.toIssue(issue)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Model(model).Association("Assignees").Replace(model.Assignees); err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Model(model).Association("Labels").Replace(model.Labels); err != nil {
		return nil, err
	}

	domainModel, err := r.Get(domain.IssueId(model.Id))
	if err != nil {
		return nil, err
	}

	if err := r.saveLabelHistory(model); err != nil {
		return nil, err
	}

	return domainModel, nil
}

func (r *IssueRepository) Update(issue *domain.Issue) (*domain.Issue, error) {
	model := r.toIssue(issue)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Model(model).Association("Assignees").Replace(model.Assignees); err != nil {
		return nil, err
	}

	if err := r.conn.GetEngine().Model(model).Association("Labels").Replace(model.Labels); err != nil {
		return nil, err
	}

	domainModel, err := r.Get(domain.IssueId(model.Id))
	if err != nil {
		return nil, err
	}

	if err := r.saveLabelHistory(model); err != nil {
		return nil, err
	}

	return domainModel, nil
}

func (r *IssueRepository) Delete(id domain.IssueId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Issue{}).Error
}

func (r *IssueRepository) toDomainIssue(issue *models.Issue) (*domain.Issue, error) {
	project, err := r.projectRepo.toDomainProject(&issue.Project)
	if err != nil {
		return nil, err
	}

	var release *domain.Release
	if issue.Release != (*models.Release)(nil) {
		release, err = r.releaseRepo.toDomainRelease(issue.Release)
		if err != nil {
			return nil, err
		}
	}

	return &domain.Issue{
		Id:        domain.IssueId(issue.Id),
		Iid:       domain.IssueIid(issue.Iid),
		Title:     issue.Title,
		IssueType: domain.IssueType(issue.IssueType),
		Assignees: utils.Map(issue.Assignees, func(a models.User) domain.User {
			return *r.userRepo.toDomainUser(&a)
		}),
		WebUrl: issue.WebUrl,
		Labels: utils.Map(issue.Labels, func(l models.Label) domain.Label {
			return *r.labelRepo.toDomainLabel(&l)
		}),
		Project: *project,
		Release: release,
		TaskType: func() *domain.Label {
			if issue.TaskType != (*models.Label)(nil) {
				return r.labelRepo.toDomainLabel(issue.TaskType)
			}
			return nil
		}(),
		EstimateDev: issue.EstimateDev,
		EstimateQA:  issue.EstimateQA,
	}, nil
}

func (r *IssueRepository) toIssue(issue *domain.Issue) *models.Issue {

	return &models.Issue{
		Id:        uint(issue.Id),
		Iid:       string(issue.Iid),
		Title:     issue.Title,
		IssueType: string(issue.IssueType),
		Assignees: utils.Map(issue.Assignees, func(a domain.User) models.User {
			return *r.userRepo.toUser(&a)
		}),
		WebUrl: issue.WebUrl,
		Labels: utils.Map(issue.Labels, func(l domain.Label) models.Label {
			return *r.labelRepo.toLabel(&l)
		}),
		ProjectId: uint(issue.Project.Id),
		ReleaseId: func() *string {
			if issue.Release != (*domain.Release)(nil) {
				val := string(issue.Release.Id)
				return &val
			}
			return nil
		}(),
		TaskTypeId: func() *string {
			if issue.TaskType != (*domain.Label)(nil) {
				val := string(issue.TaskType.Id)
				return &val
			}
			return nil
		}(),
		EstimateDev: issue.EstimateDev,
		EstimateQA:  issue.EstimateQA,
	}
}

func (r *IssueRepository) getQuery() *gorm.DB {
	query := r.conn.GetEngine().Model(&models.Issue{})
	query = query.Preload("Assignees").Preload("Labels")
	query = query.Preload("Project").Preload("Project.Team")
	query = query.Preload("Release").Preload("TaskType")
	return query
}

func (r *IssueRepository) saveLabelHistory(issue *models.Issue) error {
	domainIssue, err := r.toDomainIssue(issue)
	if err != nil {
		return err
	}

	domainHistory := make([]domain.LabelHistory, 0)
	var lastHistory models.LabelHistory
	err = r.conn.GetEngine().Where("issue_id = ?", issue.Id).Order("created_at DESC").First(&lastHistory).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		domainHistory = append(domainHistory, *r.toDomainLabelHistory(&lastHistory))
	}

	domainIssue.SetHistory(domainHistory)
	addedHistory := domainIssue.GetAddedHistory()
	for _, historyItem := range addedHistory {
		err := r.conn.GetEngine().Save(&models.LabelHistory{
			IssueId: issue.Id,
			Labels: datatypes.NewJSONSlice(utils.Map(historyItem.Labels, func(id domain.LabelId) string {
				return string(id)
			})),
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *IssueRepository) toDomainLabelHistory(history *models.LabelHistory) *domain.LabelHistory {
	return &domain.LabelHistory{
		Labels: utils.Map(history.Labels, func(id string) domain.LabelId {
			return domain.LabelId(id)
		}),
		CreatedAt: history.CreatedAt,
	}
}
