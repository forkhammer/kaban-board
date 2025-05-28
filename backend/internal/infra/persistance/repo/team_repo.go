package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/pkg/utils"

	"gorm.io/gorm"
)

type TeamRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *TeamRepository) Get(id domain.TeamId) (*domain.Team, error) {
	team := &models.Team{}
	if err := r.getQuery().Where("id = ?", id).First(team).Error; err != nil {
		return nil, err
	}
	return r.toDomainTeam(team), nil
}

func (r *TeamRepository) GetByUsername(username string) (*domain.Team, error) {
	team := &models.Team{}
	if err := r.getQuery().Where("username = ?", username).First(team).Error; err != nil {
		return nil, err
	}
	return r.toDomainTeam(team), nil
}

func (r *TeamRepository) List(spec repo.QuerySpec) (*[]domain.Team, error) {
	teams := make([]models.Team, 0)
	query := r.getQuery()

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	domainTeams := make([]domain.Team, len(teams))
	for i, team := range teams {
		domainTeams[i] = *r.toDomainTeam(&team)
	}

	return &domainTeams, nil
}

func (r *TeamRepository) Create(team *domain.Team) error {
	model := r.toTeam(team)
	return r.conn.GetEngine().Create(model).Error
}

func (r *TeamRepository) Update(team *domain.Team) error {
	model := r.toTeam(team)
	return r.conn.GetEngine().Save(model).Error
}

func (r *TeamRepository) Delete(id domain.TeamId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Team{}).Error
}

func (r *TeamRepository) toDomainTeam(team *models.Team) *domain.Team {
	return &domain.Team{
		Id:    domain.TeamId(team.Id),
		Title: team.Title,
		Groups: utils.Map(team.Groups, func(g *models.Group) domain.GroupId {
			return domain.GroupId(g.Id)
		}),
	}
}

func (r *TeamRepository) toTeam(team *domain.Team) *models.Team {
	return &models.Team{
		Id:    uint(team.Id),
		Title: team.Title,
		Groups: utils.Map(team.Groups, func(groupId domain.GroupId) *models.Group {
			return &models.Group{Id: uint(groupId)}
		}),
	}
}

func (r *TeamRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Team{}).Preload("Groups")
}
