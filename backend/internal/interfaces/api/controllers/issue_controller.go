package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IssueController struct {
	issueUC *usecases.IssueUseCases `di.inject:"IssueUseCases"`
}

func (c *IssueController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/issue", c.getIssues)
	router.GET("/issue/:id", c.getIssue)
	router.PUT("/issue/:id", c.saveIssue)
	router.POST("/issue/:id/bind", c.bindIssue)
	return nil
}

func (c *IssueController) getIssues(ctx *gin.Context) {
	var request dto.IssuesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issues, err := c.issueUC.GetIssues(&queries.IssueFilter{
		AssigneeId: request.Assignee,
		TeamId:     request.Team,
		SprintId:   request.Sprint,
		GroupId:    request.Group,
		ProjectId:  request.Project,
		Search:     request.Search,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssues(*issues))
}

func (c *IssueController) getIssue(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issue, err := c.issueUC.GetIssue(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssue(issue))
}

func (c *IssueController) bindIssue(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.BindIssueRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issue, err := c.issueUC.BindIssue(uint(id), request.SprintId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssue(issue))
}

func (c *IssueController) saveIssue(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.SaveIssueRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issue, err := c.issueUC.SaveIssue(usecases.SaveIssueRequest{
		Id:          uint(id),
		BindingId:   request.BindingId,
		EstimateDev: request.EstimateDev,
		EstimateQA:  request.EstimateQA,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssue(issue))
}
