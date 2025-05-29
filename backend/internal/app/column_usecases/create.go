package column_usecases

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"
)

type CreateColumnRequest struct {
	Name   string
	Labels []string
	TeamId *uint
}

type CreateColumnUseCase struct {
	columnRepo repo.ColumnRepo `di.inject:"ColumnRepository"`
	teamRepo   repo.TeamRepo   `di.inject:"TeamRepository"`
}

func (uc *CreateColumnUseCase) Execute(request *CreateColumnRequest) (*domain.Column, error) {
	var team *domain.Team
	var err error
	if request.TeamId != nil {
		team, err = uc.teamRepo.Get(domain.TeamId(*request.TeamId))
		if err != nil {
			return nil, err
		}
	}
	order := 0

	column, err := domain.NewColumn(
		0,
		request.Name,
		utils.Map(request.Labels, func(label string) domain.LabelId {
			return domain.LabelId(label)
		}),
		team,
		&order,
	)

	if err := column.Validate(); err != nil {
		return nil, err
	}

	column, err = uc.columnRepo.Create(column)
	if err != nil {
		return nil, err
	}

	return column, nil
}
