package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	groupUC *usecases.GroupUseCases `di.inject:"GroupUseCases"`
}

func (c *GroupController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/groups", c.getGroups)
	router.GET("/groups/:id", c.getGroupById)
	return nil
}

func (c *GroupController) getGroups(ctx *gin.Context) {
	groups, err := c.groupUC.GetGroups()

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeGroups(groups))
}

func (c *GroupController) getGroupById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	group, err := c.groupUC.GetGroup(uint(id))

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeGroup(group))
}
