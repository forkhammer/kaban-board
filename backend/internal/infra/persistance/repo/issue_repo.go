package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/gorm"
)

type IssueRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *IssueRepository) Get(id domain.IssueId) (*domain.Issue, error) {
	issue := &models.Issue{}
	if err := r.getQuery().Where("id = ?", id).First(issue).Error; err != nil {
		return nil, err
	}
	return r.toDomainIssue(issue), nil
}

func (r *IssueRepository) GetByIssuename(issuename string) (*domain.Issue, error) {
	issue := &models.Issue{}
	if err := r.getQuery().Where("issuename = ?", issuename).First(issue).Error; err != nil {
		return nil, err
	}
	return r.toDomainIssue(issue), nil
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
		domainIssues[i] = *r.toDomainIssue(&issue)
	}

	return &domainIssues, nil
}

func (r *IssueRepository) Create(issue *domain.Issue) error {
	model := r.toIssue(issue)
	return r.conn.GetEngine().Create(model).Error
}

func (r *IssueRepository) Update(issue *domain.Issue) error {
	model := r.toIssue(issue)
	return r.conn.GetEngine().Save(model).Error
}

func (r *IssueRepository) Delete(id domain.IssueId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Issue{}).Error
}

func (r *IssueRepository) toDomainIssue(issue *models.Issue) *domain.Issue {
	return &domain.Issue{
		Id:        domain.IssueId(issue.Id),
		Iid:       domain.IssueIid(issue.Iid),
		Title:     issue.Title,
		IssueType: domain.IssueType(issue.IssueType),
		Assignees: utils.Map(issue.Assignees, func(a models.User) domain.UserId {
			return domain.UserId(a.Id)
		}),
		WebUrl: issue.WebUrl,
		Labels: utils.Map(issue.Labels, func(l models.Label) domain.LabelId {
			return domain.LabelId(l.Id)
		}),
		ProjectId:   domain.ProjectId(issue.ProjectId),
		ReleaseId:   (*domain.ReleaseId)(issue.ReleaseId),
		TaskType:    (*domain.LabelId)(issue.TaskTypeId),
		EstimateDev: issue.EstimateDev,
		EstimateQA:  issue.EstimateQA,
	}
}

func (r *IssueRepository) toIssue(issue *domain.Issue) *models.Issue {
	return &models.Issue{
		Id:        string(issue.Id),
		Iid:       string(issue.Iid),
		Title:     issue.Title,
		IssueType: string(issue.IssueType),
		Assignees: utils.Map(issue.Assignees, func(a domain.UserId) models.User {
			return models.User{Id: uint(a)}
		}),
		WebUrl: issue.WebUrl,
		Labels: utils.Map(issue.Labels, func(l domain.LabelId) models.Label {
			return models.Label{Id: string(l)}
		}),
		ProjectId:   uint(issue.ProjectId),
		ReleaseId:   (*string)(issue.ReleaseId),
		TaskTypeId:  (*string)(issue.TaskType),
		EstimateDev: issue.EstimateDev,
		EstimateQA:  issue.EstimateQA,
	}
}

func (r *IssueRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Issue{}).Preload("Assignees").Preload("Labels").Preload("Project").Preload("Release").Preload("TaskType")
}
