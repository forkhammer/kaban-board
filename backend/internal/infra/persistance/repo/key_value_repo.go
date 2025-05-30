package repo

import (
	"encoding/json"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/datatypes"
)

type KeyValueRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *KeyValueRepository) Get(key string, to any, def any) error {
	defData, err := json.Marshal(def)

	if err != nil {
		return nil
	}

	defValue := models.KVElement{
		Key:   key,
		Value: datatypes.JSON(defData),
	}

	var kv models.KVElement

	if result := r.conn.GetEngine().Model(&models.KVElement{}).Where("key = ?", key).Attrs(&defValue).FirstOrCreate(&kv); result.Error != nil {
		return result.Error
	}

	if err := json.Unmarshal([]byte(kv.Value.String()), to); err != nil {
		return err
	}
	return nil
}

func (r *KeyValueRepository) Set(key string, value any) error {
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	var existValue any

	if err := r.Get(key, &existValue, value); err != nil {
		return err
	}

	kv := models.KVElement{Key: key, Value: datatypes.JSON(data)}

	if err := r.conn.GetEngine().Save(&kv).Error; err != nil {
		return err
	}

	return nil
}
