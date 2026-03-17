package controllers

import (
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/infra/hub"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SprintWSController struct {
	hub       *hub.SprintHub            `di.inject:"SprintHub"`
	accountUC *usecases.AccountUseCases `di.inject:"AccountUseCases"`
}

func (c *SprintWSController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/sprint/:id/ws", c.handleWS)
	return nil
}

func (c *SprintWSController) handleWS(ctx *gin.Context) {
	sprintId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	token := ctx.Query("token")
	_, err = c.accountUC.GetActiveUser(token)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	id := domain.SprintId(sprintId)
	c.hub.Register(id, conn)
	defer c.hub.Unregister(id, conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
