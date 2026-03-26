package usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type UpdateLabelRequest struct {
	Title         string
	AltName       *string
	BindingStatus *domain.IssueBindingStatus
	Priority      *domain.IssueBindingPriority
}

type KanbanLabel struct {
	Id            string
	Name          string
	AltName       *string
	BindingStatus *domain.IssueBindingStatus
	Priority      *domain.IssueBindingPriority
}

type LabelUseCases struct {
	labelRepo  repo.LabelRepo     `di.inject:"LabelRepository"`
	labelQuery queries.LabelQuery `di.inject:"LabelQuery"`
	groupRepo  repo.GroupRepo     `di.inject:"GroupRepository"`
	groupQuery queries.GroupQuery `di.inject:"GroupQuery"`
}

func (uc *LabelUseCases) GetLabels() ([]KanbanLabel, error) {
	labels, err := uc.labelRepo.List(nil)
	if err != nil {
		return nil, err
	}
	kanbanLabels := utils.Map(
		utils.Unique(labels, func(l domain.Label) string {
			return l.Name
		}),
		func(l domain.Label) KanbanLabel {
			return KanbanLabel{
				Id:            l.Name,
				Name:          l.Name,
				AltName:       l.AltName,
				BindingStatus: l.BindingStatus,
				Priority:      l.Priority,
			}
		},
	)
	return kanbanLabels, nil
}

func (uc *LabelUseCases) GetLabel(id string) (*domain.Label, error) {
	return uc.labelRepo.Get(domain.LabelId(id))
}

func (uc *LabelUseCases) Update(request *UpdateLabelRequest) error {
	labels, err := uc.labelRepo.List(uc.labelQuery.GetSpec(queries.LabelFilter{Names: []string{request.Title}}))
	if err != nil {
		return err
	}

	for _, label := range labels {
		label.AltName = request.AltName
		label.BindingStatus = request.BindingStatus
		label.Priority = request.Priority
		if err := label.Validate(); err != nil {
			return err
		}

		_, err := uc.labelRepo.Update(&label)
		if err != nil {
			return err
		}
	}

	return nil
}
