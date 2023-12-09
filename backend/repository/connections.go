package repository

import (
	"errors"
	"main/repository/rdb"
	"main/tools"
)

func GetRepositoryFactory(repositoryType tools.RepositoryType, connection tools.ConnectionInterface) (tools.RepositoryFactory, error) {
	switch repositoryType {
	case tools.Postgresql, tools.Mysql, tools.Sqlite:
		return getRBDRepositoryFactory(connection), nil
	}

	return nil, errors.New("Invalid repository type")
}

func getRBDRepositoryFactory(connection tools.ConnectionInterface) tools.RepositoryFactory {
	return rdb.NewRDBRepositoryFactory(connection)
}
