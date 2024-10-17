package rdb

import (
	"main/repository"
	"main/repository/models"

	"gorm.io/gorm"
)

type TeamRepository struct {
	connection repository.ConnectionInterface
}

func (r *TeamRepository) GetTeams(to *[]models.Team) error {
	result := r.getDb().Preload("Groups").Find(to)
	return result.Error
}

func (r *TeamRepository) GetTeamById(to *models.Team, id int) error {
	if result := r.getDb().Where("id = ?", id).Preload("Groups").First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) SaveTeam(team *models.Team) error {
	if result := r.getDb().Save(team); result.Error != nil {
		return result.Error
	}

	if err := r.getDb().Model(team).Association("Groups").Replace(team.Groups); err != nil {
		return err
	}

	return nil
}

func (r *TeamRepository) CreateTeam(team *models.Team) error {
	if result := r.getDb().Create(team); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) DeleteTeam(team *models.Team) error {
	result := r.getDb().Delete(team)
	return result.Error
}

func (r *TeamRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
