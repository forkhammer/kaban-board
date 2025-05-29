package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *ProjectRepository) Get(id domain.ProjectId) (*domain.Project, error) {
	project := &models.Project{}
	if err := r.conn.GetEngine().Where("id = ?", id).First(project).Error; err != nil {
		return nil, err
	}
	model, err := r.toDomainProject(project)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *ProjectRepository) List(spec repo.QuerySpec) (*[]domain.Project, error) {
	projects := make([]models.Project, 0)
	query := r.conn.GetEngine().Model(&models.Project{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}

	domainProjects := make([]domain.Project, len(projects))
	for i, Project := range projects {
		if elem, err := r.toDomainProject(&Project); err != nil {
			return nil, err
		} else {
			domainProjects[i] = *elem
		}
	}

	return &domainProjects, nil
}

func (r *ProjectRepository) Create(Project *domain.Project) (*domain.Project, error) {
	model, err := r.toProject(Project)
	if err != nil {
		return nil, err
	}

	err = r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.ProjectId(model.Id))

}

func (r *ProjectRepository) Update(Project *domain.Project) (*domain.Project, error) {
	model, err := r.toProject(Project)
	if err != nil {
		return nil, err
	}

	err = r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.ProjectId(model.Id))
}

func (r *ProjectRepository) Delete(id domain.ProjectId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Project{}).Error
}

func (r *ProjectRepository) toDomainProject(project *models.Project) (*domain.Project, error) {
	teamId, err := utils.IntToUintPtr(project.TeamId)
	if err != nil {
		return nil, err
	}

	return &domain.Project{
		Id:        domain.ProjectId(project.Id),
		Name:      project.Name,
		TeamId:    (*domain.TeamId)(teamId),
		IsVisible: project.IsVisible,
		Users: utils.Map(project.Users, func(userId int64) domain.UserId {
			return domain.UserId(userId)
		}),
	}, nil
}

func (r *ProjectRepository) toProject(column *domain.Project) (*models.Project, error) {
	teamId, err := utils.UintToIntPtr((*uint)(column.TeamId))
	if err != nil {
		return nil, nil
	}

	return &models.Project{
		Id:        uint(column.Id),
		Name:      column.Name,
		TeamId:    teamId,
		IsVisible: column.IsVisible,
		Users: datatypes.NewJSONSlice(utils.Map(column.Users, func(userId domain.UserId) int64 {
			return int64(userId)
		})),
	}, nil
}
