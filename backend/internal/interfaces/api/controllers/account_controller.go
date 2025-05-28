package controllers

import (
	"main/internal/app/account_usecases"
	domain "main/internal/domain/models"
	"main/internal/infra/persistance/models"
	"main/internal/interfaces/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	login    *account_usecases.LoginUseCase      `di.inject:"LoginUseCase"`
	active   *account_usecases.ActiveUserUseCase `di.inject:"ActiveUserUseCase"`
	register *account_usecases.RegisterUseCase   `di.inject:"RegisterUseCase"`
}

func (c *AccountController) RegisterRoutes(router *gin.Engine) error {
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

	token, err := c.login.Execute(request.Username, request.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	account, err := c.register.Execute(request.Username, request.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeAccount(account))
}
