package kanban

type KanbanSettings struct {
	TypeTaskLabels []string `json:"typeTaskLabels"`
	kvstore        *KVStore `di.inject:"kvStore"`
}

func (s *KanbanSettings) PostConstruct() error {
	return s.Init()
}

func (s *KanbanSettings) Init() error {
	if err := s.kvstore.GetValue("type_task_labels", &s.TypeTaskLabels, []string{}); err != nil {
		return err
	}

	return nil
}
