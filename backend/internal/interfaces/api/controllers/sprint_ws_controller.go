package controllers

import (
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SprintWSController struct {
	hub interfaces.SprintHub `di.inject:"SprintHub"`
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
