package kanban

import (
	"main/gitlab"
	"main/tools"
)

type LabelService struct {
	labelRepository tools.LabelRepositoryInterface `di.inject:"labelRepository"`
}

func (s *LabelService) SaveLabels(labels []gitlab.GitlabLabel) ([]Label, error) {
	var resultLabels = make([]Label, 0)

	for i := range labels {
		l := &labels[i]

		var label Label
		err := s.labelRepository.GetOrCreate(&label, Label{Id: l.Id}, Label{
			Id:        l.Id,
			Name:      l.Title,
			Color:     l.Color,
			TextColor: l.TextColor,
		})

		if err != nil {
			return resultLabels, err
		}

		resultLabels = append(resultLabels, label)
	}

	return resultLabels, nil
}

func (s *LabelService) GetAllKanbanLabels() ([]*KanbanLabel, error) {
	var labels []Label
	err := s.labelRepository.GetLabels(&labels)
	if err != nil {
		return nil, err
	}
	titles := tools.Map(labels, func(label Label) string {
		return label.Name
	})
	titles = tools.Unique(titles, func(t string) string {
		return t
	})
	kanbanLabels := tools.Map(titles, func(t string) *KanbanLabel {
		return &KanbanLabel{
			Id:    t,
			Title: t,
		}
	})
	return kanbanLabels, nil
}
