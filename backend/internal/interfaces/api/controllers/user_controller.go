package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	"main/internal/domain/repo"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepo repo.UserRepo          `di.inject:"UserRepository"`
	userUC   *usecases.UserUseCases `di.inject:"UserUseCases"`
}

func (c *UserController) RegisterRoutes(router *gin.Engine) error {
	userRoutes := router.Group("/")
	userRoutes.Use(middleware.AuthRequiredMiddleware())
	userRoutes.GET("/users", c.getUsers)
	userRoutes.GET("/users/:id", c.getUser)
	userRoutes.POST("/users/:id/visibility", c.setUserVisibility)
	userRoutes.POST("/users/:id/groups", c.setUserGroups)
	return nil
}

func (c *UserController) getUsers(ctx *gin.Context) {
	var request dto.GetUsersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	users, err := c.userUC.GetUsers(&queries.UserFilter{
		TeamId: request.TeamId,
		Search: request.Search,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeUsers(*users))
}

func (c *UserController) getUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := c.userUC.GetUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeUser(user))
}

func (c *UserController) setUserVisibility(ctx *gin.Context) {
	var request dto.SetUserVisibilityRequest

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err = ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := c.userUC.SetVisibility(uint(id), request.Visible)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeUser(user))
}

func (c *UserController) setUserGroups(ctx *gin.Context) {
	var request dto.SetUserGroupsRequest

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err = ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := c.userUC.SetGroups(uint(id), request.Groups)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeUser(user))
}
