package kanban

type KanbanSettings struct {
	TaskTypeLabels []string `json:"taskTypeLabels"`
	kvstore        *KVStore `di.inject:"kvStore"`
}

type SaveTaskTypeLabelsRequest struct {
	Labels []string `json:"labels"`
}

func (s *KanbanSettings) PostConstruct() error {
	return s.Init()
}

func (s *KanbanSettings) Init() error {
	if err := s.kvstore.GetValue("task_type_labels", &s.TaskTypeLabels, []string{}); err != nil {
		return err
	}

	return nil
}

func (s *KanbanSettings) SetTaskTypeLabels(labels []string) error {
	s.TaskTypeLabels = labels
	return s.kvstore.SetValue("task_type_labels", s.TaskTypeLabels)
}
