package repo

type KeyValueRepo interface {
	Get(key string, to any, def any) error
	Set(key string, value any) error
}
