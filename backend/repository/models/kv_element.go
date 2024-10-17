package models

import "gorm.io/datatypes"

type KVElement struct {
	Key   string         `gorm:"key;primaryKey"`
	Value datatypes.JSON `gorm:"value"`
}
