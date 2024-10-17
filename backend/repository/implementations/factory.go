package implementations

import (
	"errors"
	"main/repository"
	"main/repository/implementations/rdb"
)

func GetRepositoryFactory(repositoryType repository.RepositoryType, connection repository.ConnectionInterface) (repository.RepositoryFactory, error) {
	switch repositoryType {
	case repository.Postgresql, repository.Mysql, repository.Sqlite:
		return getRBDRepositoryFactory(connection), nil
	}

	return nil, errors.New("Invalid repository type")
}

func getRBDRepositoryFactory(connection repository.ConnectionInterface) repository.RepositoryFactory {
	return rdb.NewRDBRepositoryFactory(connection)
}
