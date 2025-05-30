package usecases

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"
)

type KanbanUseCases struct {
	userRepo  repo.UserRepo  `di.inject:"UserRepository"`
	issueRepo repo.IssueRepo `di.inject:"IssueRepository"`
}

func (uc *KanbanUseCases) GetBoard() (*models.Board, error) {
	users, err := uc.userRepo.List(nil)
	if err != nil {
		return nil, err
	}

	issues, err := uc.issueRepo.List(nil)
	if err != nil {
		return nil, err
	}

	kanban, err := models.NewKanban(*users, *issues)
	if err != nil {
		return nil, err
	}

	return kanban.GetBoard()
}
