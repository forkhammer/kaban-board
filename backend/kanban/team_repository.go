package kanban

import "main/db"

type TeamRepository struct{}

func (r *TeamRepository) GetTeams() ([]Team, error) {
	var teams []Team
	result := db.DefaultConnection.Db.Find(&teams)
	return teams, result.Error
}

func (r *TeamRepository) GetTeamById(id int) (*Team, error) {
	var team Team
	if result := db.DefaultConnection.Db.Where(&Team{Id: id}).First(&team); result.Error != nil {
		return nil, result.Error
	} else {
		return &team, nil
	}
}

func (r *TeamRepository) SaveTeam(team *Team) error {
	if result := db.DefaultConnection.Db.Save(team); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) CreateTeam(team *Team) error {
	if result := db.DefaultConnection.Db.Create(team); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *TeamRepository) DeleteTeam(id int) error {
	result := db.DefaultConnection.Db.Delete(&Team{Id: id})
	return result.Error
}
