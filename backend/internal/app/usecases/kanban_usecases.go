package usecases

import (
	"main/internal/app/queries"
	"main/internal/domain/models"
	"main/internal/domain/repo"
)

type KanbanUseCases struct {
	userRepo  repo.UserRepo     `di.inject:"UserRepository"`
	issueRepo repo.IssueRepo    `di.inject:"IssueRepository"`
	userQuery queries.UserQuery `di.inject:"UserQuery"`
}

func (uc *KanbanUseCases) GetBoard() (*models.Board, error) {
	users, err := uc.userRepo.List(uc.userQuery.OnlyVisible())
	if err != nil {
		return nil, err
	}

	issues, err := uc.issueRepo.List(nil)
	if err != nil {
		return nil, err
	}

	kanban, err := models.NewKanban(users, issues)
	if err != nil {
		return nil, err
	}

	return kanban.GetBoard()
}
