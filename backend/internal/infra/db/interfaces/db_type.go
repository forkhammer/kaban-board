package interfaces

type DbType string

const (
	Postgresql DbType = "postgresql"
	Mysql      DbType = "mysql"
	Sqlite     DbType = "sqlite"
)
