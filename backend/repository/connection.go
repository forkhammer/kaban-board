package repository

type ConnectionInterface interface {
	GetEngine() interface{}
	Migrate(models ...interface{}) error
}
