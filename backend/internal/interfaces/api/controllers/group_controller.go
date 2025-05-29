package controllers

import (
	"main/internal/app/group_usecases"
	"main/internal/interfaces/api/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	groupUC *group_usecases.GroupUseCases `di.inject:"GroupUseCases"`
}

func (c *GroupController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/groups", c.getGroups)
	router.GET("/groups/:id", c.getGroupById)
	return nil
}

func (c *GroupController) getGroups(ctx *gin.Context) {
	groups, err := c.groupUC.GetGroups()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeGroups(*groups))
}

func (c *GroupController) getGroupById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	group, err := c.groupUC.GetGroup(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeGroup(group))
}
