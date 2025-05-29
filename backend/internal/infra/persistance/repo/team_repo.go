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
	conn      interfaces.ConnectionInterface `di.inject:"db"`
	groupRepo *GroupRepository               `di.inject:"GroupRepository"`
}

func (r *TeamRepository) Get(id domain.TeamId) (*domain.Team, error) {
	team := &models.Team{}
	if err := r.getQuery().Where("id = ?", id).First(team).Error; err != nil {
		return nil, err
	}
	return r.ToDomainTeam(team)
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
		team, err := r.ToDomainTeam(&team)
		if err != nil {
			return nil, err
		}
		domainTeams[i] = *team
	}

	return &domainTeams, nil
}

func (r *TeamRepository) Create(team *domain.Team) (*domain.Team, error) {
	model := r.ToTeam(team)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}

	err = r.conn.GetEngine().Model(model).Association("Groups").Replace(model.Groups)
	if err != nil {
		return nil, err
	}
	return r.Get(domain.TeamId(model.Id))
}

func (r *TeamRepository) Update(team *domain.Team) (*domain.Team, error) {
	model := r.ToTeam(team)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}

	err = r.conn.GetEngine().Model(model).Association("Groups").Replace(model.Groups)
	if err != nil {
		return nil, err
	}

	return r.Get(domain.TeamId(model.Id))
}

func (r *TeamRepository) Delete(id domain.TeamId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Team{}).Error
}

func (r *TeamRepository) ToDomainTeam(team *models.Team) (*domain.Team, error) {
	return domain.NewTeam(
		domain.TeamId(team.Id),
		team.Title,
		utils.Map(team.Groups, func(group *models.Group) domain.Group {
			return *r.groupRepo.toDomainGroup(group)
		}),
	)
}

func (r *TeamRepository) ToTeam(team *domain.Team) *models.Team {
	return &models.Team{
		Id:    uint(team.Id),
		Title: team.Title,
		Groups: utils.Map(team.Groups, func(group domain.Group) *models.Group {
			return r.groupRepo.toGroup(&group)
		}),
	}
}

func (r *TeamRepository) getQuery() *gorm.DB {
	return r.conn.GetEngine().Model(&models.Team{}).Preload("Groups")
}
