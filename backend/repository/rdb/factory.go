package rdb

import (
	"main/tools"
)

type RDBRepositoryFactory struct {
	connection tools.ConnectionInterface
}

func NewRDBRepositoryFactory(connection tools.ConnectionInterface) *RDBRepositoryFactory {
	return &RDBRepositoryFactory{
		connection: connection,
	}
}

func (f *RDBRepositoryFactory) GetAccountRepository() tools.AccountRepositoryInterface {
	return &AccountRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetColumnRepository() tools.ColumnRepositoryInterface {
	return &ColumnRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetLabelRepository() tools.LabelRepositoryInterface {
	return &LabelRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetProjectRepository() tools.ProjectRepositoryInterface {
	return &ProjectRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetTeamRepository() tools.TeamRepositoryInterface {
	return &TeamRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetUserRepository() tools.UserRepositoryInterface {
	return &UserRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetKVStoreRepository() tools.KVStoreRepositoryInterface {
	return &KVStoreRepository{
		connection: f.connection,
	}
}

func (f *RDBRepositoryFactory) GetGroupRepository() tools.GroupRepositoryInterface {
	return &GroupRepository{
		connection: f.connection,
	}
}
