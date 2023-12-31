package rdb

import (
	"main/tools"

	"gorm.io/gorm"
)

type TeamRepository struct {
	connection tools.ConnectionInterface
}

func (r *TeamRepository) GetTeams(to interface{}) error {
	result := r.getDb().Find(to)
	return result.Error
}

func (r *TeamRepository) GetTeamById(to interface{}, id int) error {
	if result := r.getDb().Where("id = ?", id).First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) SaveTeam(team interface{}) error {
	if result := r.getDb().Save(team); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) CreateTeam(team interface{}) error {
	if result := r.getDb().Create(team); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) DeleteTeam(id int) error {
	result := r.getDb().Delete("id = ?", id)
	return result.Error
}

func (r *TeamRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
