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
	conn     interfaces.ConnectionInterface `di.inject:"db"`
	teamRepo *TeamRepository                `di.inject:"TeamRepository"`
}

func (r *ProjectRepository) Get(id domain.ProjectId) (*domain.Project, error) {
	project := &models.Project{}
	if err := r.getQuery().Where("id = ?", id).First(project).Error; err != nil {
		return nil, err
	}
	model, err := r.toDomainProject(project)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *ProjectRepository) List(spec repo.QuerySpec) ([]domain.Project, error) {
	projects := make([]models.Project, 0)
	query := r.getQuery()

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

	return domainProjects, nil
}

func (r *ProjectRepository) Count(spec repo.QuerySpec) (int64, error) {
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

	return count, nil
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
	var team *domain.Team

	if project.Team != (*models.Team)(nil) {
		var err error
		if team, err = r.teamRepo.ToDomainTeam(project.Team); err != nil {
			return nil, err
		}
	}

	return &domain.Project{
		Id:        domain.ProjectId(project.Id),
		Name:      project.Name,
		Team:      team,
		IsVisible: project.IsVisible,
		Users: utils.Map(project.Users, func(userId int64) domain.UserId {
			return domain.UserId(userId)
		}),
	}, nil
}

func (r *ProjectRepository) toProject(project *domain.Project) (*models.Project, error) {
	var teamId *uint
	if project.Team != (*domain.Team)(nil) {
		val := uint(project.Team.Id)
		teamId = &val
	}

	return &models.Project{
		Id:        uint(project.Id),
		Name:      project.Name,
		TeamId:    teamId,
		IsVisible: project.IsVisible,
		Users: datatypes.NewJSONSlice(utils.Map(project.Users, func(userId domain.UserId) int64 {
			return int64(userId)
		})),
	}, nil
}

func (r *ProjectRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Project{}).Preload("Team")
}
