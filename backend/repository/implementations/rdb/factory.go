package rdb

import (
	"main/repository"
)

type RDBRepositoryFactory struct {
	connection repository.ConnectionInterface
}

func NewRDBRepositoryFactory(connection repository.ConnectionInterface) *RDBRepositoryFactory {
	return &RDBRepositoryFactory{
		connection: connection,
	}
}

func (f *RDBRepositoryFactory) GetAccountRepository() repository.AccountRepositoryInterface {
	return &AccountRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetColumnRepository() repository.ColumnRepositoryInterface {
	return &ColumnRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetLabelRepository() repository.LabelRepositoryInterface {
	return &LabelRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetProjectRepository() repository.ProjectRepositoryInterface {
	return &ProjectRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetTeamRepository() repository.TeamRepositoryInterface {
	return &TeamRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetUserRepository() repository.UserRepositoryInterface {
	return &UserRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetKVStoreRepository() repository.KVStoreRepositoryInterface {
	return &KVStoreRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetGroupRepository() repository.GroupRepositoryInterface {
	return &GroupRepository{
		connection: f.connection,
	}
}
