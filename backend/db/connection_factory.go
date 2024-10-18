package db

import (
	"errors"
	"main/config"
	"main/repository"
)

type RDBConnectionConfig struct {
	host   string
	port   int
	dbName string
	user   string
	pass   string
}

func GetConnectionByType(repositoryType repository.RepositoryType, settings *config.Config) (repository.ConnectionInterface, error) {
	switch repositoryType {
	case repository.Postgresql:
		return NewPostgresqlConnection(
			settings.PostgresHost,
			settings.PostgresPort,
			settings.PostgresDb,
			settings.PostgresUser,
			settings.PostgresPass,
		)
	case repository.Mysql:
		return NewMysqlConnection(
			settings.MysqlHost,
			settings.MysqlPort,
			settings.MysqlDb,
			settings.MysqlUser,
			settings.MysqlPass,
		)
	case repository.Sqlite:
		return NewSqliteConnection(settings.SqliteDbFile)
	}

	return nil, errors.New("Invalid respository type")
}
