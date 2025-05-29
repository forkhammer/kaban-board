package column_usecases

import (
	"fmt"
	app_services "main/internal/app/services"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type UpdateColumnRequest struct {
	Id     uint
	Name   string
	Labels []string
	TeamId *uint
}

type UpdateColumnUseCase struct {
	columnRepo   repo.ColumnRepo            `di.inject:"ColumnRepository"`
	teamRepo     repo.TeamRepo              `di.inject:"TeamRepository"`
	labelService *app_services.LabelService `di.inject:"LabelService"`
}

func (uc *UpdateColumnUseCase) Execute(request *UpdateColumnRequest) (*domain.Column, error) {
	existLabels, diff, err := uc.labelService.ExistNames(request.Labels)
	if err != nil {
		return nil, err
	}
	if !existLabels {
		return nil, fmt.Errorf("Labels not found: %v", diff)
	}

	column, err := uc.columnRepo.Get(domain.ColumnId(request.Id))
	if err != nil {
		return nil, err
	}
	column.Name = request.Name
	column.Labels = utils.Map(request.Labels, func(label string) domain.LabelId {
		return domain.LabelId(label)
	})

	var team *domain.Team
	if request.TeamId != nil {
		team, err = uc.teamRepo.Get(domain.TeamId(*request.TeamId))
		if err != nil {
			return nil, err
		}
	}
	column.Team = team

	if err := column.Validate(); err != nil {
		return nil, err
	}

	column, err = uc.columnRepo.Update(column)
	if err != nil {
		return nil, err
	}

	return column, nil
}
