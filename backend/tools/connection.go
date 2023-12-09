package tools

type ConnectionInterface interface {
	GetEngine() interface{}
	Migrate(models ...interface{}) error
}
