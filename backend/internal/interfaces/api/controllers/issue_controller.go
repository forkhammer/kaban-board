package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"main/internal/interfaces/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IssueController struct {
	issueUC *usecases.IssueUseCases `di.inject:"IssueUseCases"`
}

func (c *IssueController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/issue", c.getIssues)
	router.GET("/issue/:id", c.getIssue)
	privateRoutes := router.Group("/")
	privateRoutes.Use(middleware.AuthRequiredMiddleware())
	privateRoutes.POST("/issue/bind", c.bindIssue)
	return nil
}

func (c *IssueController) getIssues(ctx *gin.Context) {
	var request dto.IssuesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issueFilter := queries.IssueFilter{
		AssigneeId: request.Assignee,
		TeamId:     request.Team,
		SprintId:   request.Sprint,
		GroupId:    request.Group,
		ProjectId:  request.Project,
		Search:     request.Search,
	}
	issuePage, err := c.issueUC.GetIssues(&issueFilter, request.Page, request.Limit)
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssuePage(issuePage))
}

func (c *IssueController) getIssue(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issue, err := c.issueUC.GetIssue(uint(id))
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssue(issue))
}

func (c *IssueController) bindIssue(ctx *gin.Context) {
	var request dto.BindIssueRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	result, err := c.issueUC.BindIssues(request.IssueIds, request.SprintId, request.AssigneeId)
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.BindIssueResponse{
		Results: dto.SerializeIssueBindings(result.Bindings, currentAccount),
		Errors:  result.Errors,
	})
}
