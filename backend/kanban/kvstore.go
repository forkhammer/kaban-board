package kanban

import (
	"encoding/json"
	"errors"
	"main/tools"

	"gorm.io/datatypes"
)

type KVElement struct {
	Key   string         `gorm:"key;primaryKey"`
	Value datatypes.JSON `gorm:"value"`
}

type KVStore struct {
	kvStoreRepository tools.KVStoreRepositoryInterface `di.inject:"kvStoreRepository"`
}

func (s *KVStore) Get(key string, def interface{}) *KVElement {
	data, err := json.Marshal(def)

	if err != nil {
		return nil
	}

	kv := &KVElement{Key: key, Value: datatypes.JSON(data)}
	if err := s.kvStoreRepository.GetOrCreate(key, &kv); err != nil {
		return nil
	}

	return kv
}

func (s *KVStore) GetValue(key string, to interface{}, def interface{}) error {
	if kv := s.Get(key, def); kv != nil {
		if err := json.Unmarshal([]byte(kv.Value.String()), to); err != nil {
			return err
		}
	}

	return nil
}

func (s *KVStore) SetValue(key string, value interface{}) error {
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	kv := s.Get(key, value)

	if kv != nil {
		kv.Value = datatypes.JSON(data)
		s.kvStoreRepository.Save(kv)
	} else {
		return errors.New("Key not found")
	}

	return nil

}
