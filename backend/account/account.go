package account

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	Id        uint           `gorm:"id;primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"username;not null;unique" json:"username"`
	Password  string         `gorm:"password;not null" json:"-"`
	Name      string         `gorm:"name" json:"name"`
	IsActive  bool           `gorm:"is_active;not null" json:"isActive"`
}
