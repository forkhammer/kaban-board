package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KanbanController struct {
	kanbanUC *usecases.KanbanUseCases `di.inject:"KanbanUseCases"`
}

func (c *KanbanController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/kanban-users", c.getKanbanUsers)
	return nil
}

func (c *KanbanController) getKanbanUsers(ctx *gin.Context) {
	board, err := c.kanbanUC.GetBoard()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.BoardDto{
		Users:      dto.SerializeKanbanUsers(board.Users),
		UpdateTime: board.UpdateTime,
	})
}
