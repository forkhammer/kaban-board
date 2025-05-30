package interfaces

import "gorm.io/gorm"

type ConnectionInterface interface {
	GetEngine() *gorm.DB
	Migrate(models ...any) error
}
