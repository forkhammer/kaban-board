package controllers

import (
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/infra/persistance/models"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountUC *usecases.AccountUseCases `di.inject:"AccountUseCases"`
}

func (c *AccountController) RegisterRoutes(router gin.IRouter) error {
	router.POST("/account/login", c.Login)
	router.GET("/account/user", c.GetActiveAccount)
	router.POST("/account/register", c.Register)
	return nil
}

func (c *AccountController) Login(ctx *gin.Context) {
	var request dto.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.accountUC.Login(request.Username, request.Password)

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

func (c *AccountController) GetActiveAccount(ctx *gin.Context) {
	account, found := ctx.Get("account")

	if found && account != (*models.Account)(nil) {
		ctx.JSON(http.StatusOK, dto.NewActiveUserResponse(account.(*domain.Account)))
	} else {
		ctx.JSON(http.StatusNotFound, dto.ErrorsResponse{Errors: []string{"User not found"}})
	}
}

func (c *AccountController) Register(ctx *gin.Context) {
	var request dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, err := c.accountUC.Register(request.Username, request.Password)

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeAccount(account))
}
