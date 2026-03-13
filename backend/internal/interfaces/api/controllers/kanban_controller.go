package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KanbanController struct {
	kanbanUC *usecases.KanbanUseCases `di.inject:"KanbanUseCases"`
}

func (c *KanbanController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/kanban-users", c.getKanbanUsers)
	return nil
}

func (c *KanbanController) getKanbanUsers(ctx *gin.Context) {
	board, err := c.kanbanUC.GetBoard()

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.BoardDto{
		Users:      dto.SerializeKanbanUsers(board.Users),
		UpdateTime: board.UpdateTime,
	})
}
