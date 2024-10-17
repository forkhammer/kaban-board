package rdb

import (
	"main/repository"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	connection repository.ConnectionInterface
}

func (r *ProjectRepository) GetProjectById(to interface{}, id uint) error {
	if result := r.getDb().Where("id = ?", id).First(to); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *ProjectRepository) GetProjects(to interface{}) error {
	result := r.getDb().Order("name").Find(to)
	return result.Error
}

func (r *ProjectRepository) SaveProject(project interface{}) error {
	result := r.getDb().Save(project)
	return result.Error
}

func (r *ProjectRepository) CreateProject(project interface{}) error {
	result := r.getDb().Create(project)
	return result.Error
}

func (r *ProjectRepository) getDb() *gorm.DB {
	return r.connection.GetEngine().(*gorm.DB)
}
