package kanban

import "main/db"

type ProjectRepository struct{}

func (r *ProjectRepository) GetProjectById(id uint) (*Project, error) {
	var project Project
	if result := db.DefaultConnection.Db.Where(&Project{Id: id}).First(&project); result.Error != nil {
		return nil, result.Error
	} else {
		return &project, nil
	}
}

func (r *ProjectRepository) GetProjects() ([]Project, error) {
	var projects []Project
	result := db.DefaultConnection.Db.Order("name").Find(&projects)
	return projects, result.Error
}

func (r *ProjectRepository) SaveProject(project *Project) error {
	result := db.DefaultConnection.Db.Save(project)
	return result.Error
}

func (r *ProjectRepository) CreateProject(project *Project) error {
	result := db.DefaultConnection.Db.Create(project)
	return result.Error
}
