package label_usecases

import (
	"main/internal/app/queries"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type UpdateLabelRequest struct {
	Title   string
	AltName *string
}

type LabelUseCases struct {
	labelRepo  repo.LabelRepo     `di.inject:"LabelRepository"`
	labelQuery queries.LabelQuery `di.inject:"LabelQuery"`
	groupRepo  repo.GroupRepo     `di.inject:"GroupRepository"`
	groupQuery queries.GroupQuery `di.inject:"GroupQuery"`
}

func (uc *LabelUseCases) GetLabels() (*[]domain.Label, error) {
	return uc.labelRepo.List(nil)
}

func (uc *LabelUseCases) GetLabel(id string) (*domain.Label, error) {
	return uc.labelRepo.Get(domain.LabelId(id))
}

func (uc *LabelUseCases) Update(request *UpdateLabelRequest) error {
	labels, err := uc.labelRepo.List(uc.labelQuery.GetSpec(queries.LabelFilter{Names: []string{request.Title}}))
	if err != nil {
		return err
	}

	for _, label := range *labels {
		label.AltName = request.AltName
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
