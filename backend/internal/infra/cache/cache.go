package cache

import "time"

type Cache interface {
	Set(key string, value any, duration time.Duration)
	Get(key string) (any, bool)
	Delete(key string) error
	GetByPrefix(prefix string) map[string]any
}
