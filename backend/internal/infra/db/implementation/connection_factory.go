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

func GetConnectionByType(dbType interfaces.DbType, cfg *config.Config) (interfaces.ConnectionInterface, error) {
	var conn interfaces.ConnectionInterface
	var err error

	switch dbType {
	case interfaces.Postgresql:
		conn, err = NewPostgresqlConnection(cfg)
	case interfaces.Mysql:
		conn, err = NewMysqlConnection(cfg)
	case interfaces.Sqlite:
		conn, err = NewSqliteConnection(cfg.SqliteDbFile, cfg)
	default:
		return nil, errors.New("Invalid repository type")
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}
