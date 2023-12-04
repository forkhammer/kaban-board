package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService AccountService
}

func (c *AccountController) RegisterRoutes(engine *gin.Engine) {
	engine.POST("/account/register", c.Register)
	engine.POST("/account/login", c.Login)
	engine.GET("/account/user", c.GetActiveAccount)
}

func (c *AccountController) Register(ctx *gin.Context) {
	var request RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := c.accountService.RegisterAccount(request.Username, request.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) Login(ctx *gin.Context) {
	var request LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.accountService.GetTokenByCredentials(request.Username, request.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *AccountController) GetActiveAccount(ctx *gin.Context) {
	account, found := ctx.Get("account")

	if found {
		ctx.JSON(http.StatusOK, gin.H{"user": account})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": []string{"User not found"}})
	}
}
