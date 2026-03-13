package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GitLabAuthController struct {
	gitLabAuthUC *usecases.GitLabAuthUseCases `di.inject:"GitLabAuthUseCases"`
}

func (c *GitLabAuthController) RegisterRoutes(router gin.IRouter) error {
	auth := router.Group("/auth/gitlab")
	auth.GET("", c.GetAuthorizationURL)
	auth.GET("/callback", c.Callback)
	return nil
}

func (c *GitLabAuthController) GetAuthorizationURL(ctx *gin.Context) {
	url, err := c.gitLabAuthUC.GetAuthorizationURL()
	if utils.HandleException(ctx, err) {
		return
	}

	if url == "" {
		ctx.JSON(http.StatusOK, dto.GitLabAuthConfigResponse{
			Enabled: false,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.GitLabAuthConfigResponse{
		Enabled: true,
		URL:     url,
	})
}

func (c *GitLabAuthController) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	token, account, err := c.gitLabAuthUC.Login(code)
	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.GitLabAuthCallbackResponse{
		Token: token,
		User:  *dto.SerializeAccount(account),
	})
}
