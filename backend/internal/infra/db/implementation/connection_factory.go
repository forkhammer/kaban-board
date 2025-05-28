package implementation

import (
	"errors"
	"main/config"
	"main/internal/infra/db/interfaces"
)

type RDBConnectionConfig struct {
	host   string
	port   int
	dbName string
	user   string
	pass   string
}

func GetConnectionByType(dbType interfaces.DbType, settings *config.Config) (interfaces.ConnectionInterface, error) {
	switch dbType {
	case interfaces.Postgresql:
		return NewPostgresqlConnection(
			settings.PostgresHost,
			settings.PostgresPort,
			settings.PostgresDb,
			settings.PostgresUser,
			settings.PostgresPass,
		)
	case interfaces.Mysql:
		return NewMysqlConnection(
			settings.MysqlHost,
			settings.MysqlPort,
			settings.MysqlDb,
			settings.MysqlUser,
			settings.MysqlPass,
		)
	case interfaces.Sqlite:
		return NewSqliteConnection(settings.SqliteDbFile)
	}

	return nil, errors.New("Invalid respository type")
}
