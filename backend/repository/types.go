package repository

import "main/repository/models"

type RepositoryType string

const (
	Postgresql RepositoryType = "postgresql"
	Mysql      RepositoryType = "mysql"
	Sqlite     RepositoryType = "sqlite"
	Redis      RepositoryType = "redis"
)

type RepositoryFactory interface {
	GetAccountRepository() AccountRepositoryInterface
	GetColumnRepository() ColumnRepositoryInterface
	GetLabelRepository() LabelRepositoryInterface
	GetProjectRepository() ProjectRepositoryInterface
	GetTeamRepository() TeamRepositoryInterface
	GetUserRepository() UserRepositoryInterface
	GetKVStoreRepository() KVStoreRepositoryInterface
	GetGroupRepository() GroupRepositoryInterface
}

type AccountRepositoryInterface interface {
	CreateAccount(account interface{}) error
	GetAccountByUsername(to interface{}, username string) error
	GetAccountById(to interface{}, id uint) error
}

type ColumnRepositoryInterface interface {
	GetColumns(to interface{}) error
	GetColumnById(to interface{}, id int) error
	SaveColumn(column interface{}) error
	CreateColumn(column interface{}) error
	DeleteColumn(column interface{}) error
}

type LabelRepositoryInterface interface {
	GetOrCreate(to interface{}, query interface{}, attrs interface{}) error
	GetLabels(to interface{}) error
	SaveLabel(label interface{}) error
	GetLabelsByName(to interface{}, title string) error
}

type ProjectRepositoryInterface interface {
	GetProjectById(to interface{}, id uint) error
	GetProjects(to interface{}) error
	SaveProject(project interface{}) error
	CreateProject(project interface{}) error
}

type TeamRepositoryInterface interface {
	GetTeams(to *[]models.Team) error
	GetTeamById(to *models.Team, id int) error
	SaveTeam(team *models.Team) error
	CreateTeam(team *models.Team) error
	DeleteTeam(team *models.Team) error
}

type UserRepositoryInterface interface {
	GetUsers(to *[]models.User) error
	GetVisibleUsers(to *[]models.User) error
	GetOrCreate(to *models.User, query, attrs interface{}) error
	GetUserBydId(to *models.User, id int) error
	SaveUser(user *models.User) error
}

type KVStoreRepositoryInterface interface {
	GetAll(to interface{}) error
	GetOrCreate(key string, to interface{}) error
	Save(value interface{}) error
	Delete(key string) error
}

type GroupRepositoryInterface interface {
	GetGroups(to interface{}) error
	GetGroupById(to interface{}, id int) error
	GetGroupsByIds(to interface{}, ids []int) error
	SaveGroup(group interface{}) error
	CreateGroup(group interface{}) error
	DeleteGroup(group interface{}) error
}
