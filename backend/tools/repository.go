package tools

type RepositoryType string

const (
	Postgresql RepositoryType = "postgresql"
	Mysql                     = "mysql"
	Sqlite                    = "sqlite"
	Redis                     = "redis"
)

type RepositoryFactory interface {
	GetAccountRepository() AccountRepositoryInterface
	GetColumnRepository() ColumnRepositoryInterface
	GetLabelRepository() LabelRepositoryInterface
	GetProjectRepository() ProjectRepositoryInterface
	GetTeamRepository() TeamRepositoryInterface
	GetUserRepository() UserRepositoryInterface
	GetKVStoreRepository() KVStoreRepositoryInterface
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
	GetTeams(to interface{}) error
	GetTeamById(to interface{}, id int) error
	SaveTeam(team interface{}) error
	CreateTeam(team interface{}) error
	DeleteTeam(team interface{}) error
}

type UserRepositoryInterface interface {
	GetUsers(to interface{}) error
	GetVisibleUsers(to interface{}) error
	GetOrCreate(to, query, attrs interface{}) error
	GetUserBydId(to interface{}, id int) error
	SaveUser(user interface{}) error
}

type KVStoreRepositoryInterface interface {
	GetAll(to interface{}) error
	GetOrCreate(key string, to interface{}) error
	Save(value interface{}) error
	Delete(key string) error
}
