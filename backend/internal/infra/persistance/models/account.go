package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id        uint `gorm:"id;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"username;not null;unique"`
	Password  string         `gorm:"password;not null"`
	Name      string         `gorm:"name"`
	IsActive  bool           `gorm:"is_active;not null"`
}
