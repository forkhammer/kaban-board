package db

import (
	"errors"
	"main/config"
	"main/tools"
)

type RDBConnectionConfig struct {
	host   string
	port   int
	dbName string
	user   string
	pass   string
}

func GetConnectionByType(repositoryType tools.RepositoryType, settings *config.Config) (tools.ConnectionInterface, error) {
	switch repositoryType {
	case tools.Postgresql:
		return NewPostgresqlConnection(
			settings.PostgresHost,
			settings.PostgresPort,
			settings.PostgresDb,
			settings.PostgresUser,
			settings.PostgresPass,
		)
	case tools.Mysql:
		return NewMysqlConnection(
			settings.MysqlHost,
			settings.MysqlPort,
			settings.MysqlDb,
			settings.MysqlUser,
			settings.MysqlPass,
		)
	}

	return nil, errors.New("Invalid respository type")
}
