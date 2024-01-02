package kanban

import (
	"encoding/json"
	"main/tools"

	"gorm.io/datatypes"
)

type KVElement struct {
	Key   string         `gorm:"key;primaryKey"`
	Value datatypes.JSON `gorm:"value"`
}

type KVStore struct {
	keys              []KVElement
	kvStoreRepository tools.KVStoreRepositoryInterface `di.inject:"kvStoreRepository"`
}

func (s *KVStore) PostConstruct() error {
	return s.Init()
}

func (s *KVStore) Init() error {
	return s.kvStoreRepository.GetAll(&s.keys)
}

func (s *KVStore) Get(key string, def interface{}) *KVElement {
	kv := tools.Find[KVElement](s.keys, func(item KVElement) bool {
		return item.Key == key
	})

	if kv == nil {
		data, err := json.Marshal(def)

		if err != nil {
			return nil
		}

		kv = &KVElement{Key: key, Value: datatypes.JSON(data)}
		if err := s.kvStoreRepository.Save(kv); err != nil {
			return nil
		}
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
