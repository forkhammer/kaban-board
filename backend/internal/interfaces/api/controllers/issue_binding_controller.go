package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IssueBindingController struct {
	bindingUC *usecases.IssueBindingUseCases `di.inject:"IssueBindingUseCases"`
}

func (c *IssueBindingController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/binding", c.getBindings)
	router.GET("/binding/:id", c.getBinding)
	privateRoutes := router.Group("/")
	privateRoutes.Use(middleware.AuthRequiredMiddleware())
	privateRoutes.DELETE("/binding/:id", c.deleteBinding)
	privateRoutes.PUT("/binding/:id", c.saveBinding)
	return nil
}

func (c *IssueBindingController) getBindings(ctx *gin.Context) {
	var request dto.IssuesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issueFilter := queries.IssueBindingFilter{
		AssigneeId: request.Assignee,
		TeamId:     request.Team,
		SprintId:   request.Sprint,
		GroupId:    request.Group,
		ProjectId:  request.Project,
		Search:     request.Search,
	}
	bindingPage, err := c.bindingUC.GetBindings(&issueFilter, request.Page, request.Limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssueBindingPage(bindingPage))
}

func (c *IssueBindingController) getBinding(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	issue, err := c.bindingUC.GetBinding(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssueBinding(issue))
}

func (c *IssueBindingController) deleteBinding(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = c.bindingUC.DeleteBinding(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *IssueBindingController) saveBinding(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.SaveIssueBindingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	binding, err := c.bindingUC.SaveBinding(usecases.SaveIssueBindingRequest{
		Id:          uint(id),
		EstimateDev: request.EstimateDev,
		EstimateQA:  request.EstimateQA,
		BindStatus:  (*domain.IssueBindingStatus)(request.BindStatus),
		Assignee:    request.Assignee,
		Comment:     request.Comment,
		Priority:    (*domain.IssueBindingPriority)(request.Priority),
		ReleaseId:   request.Release,
		EpicId:      request.Epic,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssueBinding(binding))
}
