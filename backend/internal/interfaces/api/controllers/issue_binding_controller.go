package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	apiutils "main/internal/interfaces/api/utils"
	"main/pkg/utils"
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
	privateRoutes.POST("/binding", c.createBinding)
	privateRoutes.POST("/binding/:id/copy", c.copyBinding)
	privateRoutes.POST("/binding/:id/move", c.moveBinding)
	privateRoutes.DELETE("/binding/:id", c.deleteBinding)
	privateRoutes.PUT("/binding/:id", c.saveBinding)
	adminRoutes := router.Group("/")
	adminRoutes.Use(middleware.AdminRequiredMiddleware())
	adminRoutes.POST("/binding/save_ordering", c.saveOrdering)
	return nil
}

func (c *IssueBindingController) getBindings(ctx *gin.Context) {
	var request dto.IssuesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	issueFilter := queries.IssueBindingFilter{
		AssigneeId: request.Assignee,
		TeamId:     request.Team,
		SprintId:   request.Sprint,
		GroupId:    request.Group,
		ProjectId:  request.Project,
		Search:     request.Search,
	}
	bindingPage, err := c.bindingUC.GetBindings(&issueFilter, request.Page, request.Limit, currentAccount)
	if apiutils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssueBindingPage(bindingPage, currentAccount))
}

func (c *IssueBindingController) getBinding(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	issue, err := c.bindingUC.GetBinding(uint(id), currentAccount)
	if apiutils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeIssueBinding(issue, currentAccount))
}

func (c *IssueBindingController) deleteBinding(ctx *gin.Context) {
	if !apiutils.IsAdminAccount(ctx) {
		ctx.Status(http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	_, err = c.bindingUC.DeleteBinding(uint(id))
	if apiutils.HandleException(ctx, err) {
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

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	binding, err := c.bindingUC.SaveBinding(usecases.SaveIssueBindingRequest{
		Id:          uint(id),
		Title:       request.Title,
		EstimateDev: request.EstimateDev,
		EstimateQA:  request.EstimateQA,
		BindStatus:  (*domain.IssueBindingStatus)(request.BindStatus),
		Assignee:    request.Assignee,
		Comment:     request.Comment,
		Priority:    (*domain.IssueBindingPriority)(request.Priority),
		ReleaseId:   request.Release,
		EpicId:      request.Epic,
		Version:     request.Version,
	}, currentAccount)
	if apiutils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeIssueBinding(binding, currentAccount))
}

func (c *IssueBindingController) saveOrdering(ctx *gin.Context) {
	request := make(dto.SetIssueBindingOrderRequest, 0)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ordering := utils.Map(request, func(o dto.SetIssueBindingOrder) usecases.IssueBindingOrdering {
		return usecases.IssueBindingOrdering{Id: o.Id, Order: o.Order}
	})

	if _, err := c.bindingUC.SaveOrdering(ordering); apiutils.HandleException(ctx, err) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *IssueBindingController) copyBinding(ctx *gin.Context) {
	if !apiutils.IsAdminAccount(ctx) {
		ctx.Status(http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.CopyIssueBindingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	binding, err := c.bindingUC.CopyBinding(usecases.CopyIssueBindingRequest{
		Id:       uint(id),
		SprintId: request.SprintId,
	}, currentAccount)
	if apiutils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, dto.SerializeIssueBinding(binding, currentAccount))
}

func (c *IssueBindingController) moveBinding(ctx *gin.Context) {
	if !apiutils.IsAdminAccount(ctx) {
		ctx.Status(http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.MoveIssueBindingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	binding, err := c.bindingUC.MoveBinding(usecases.MoveIssueBindingRequest{
		Id:       uint(id),
		SprintId: request.SprintId,
	}, currentAccount)
	if apiutils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeIssueBinding(binding, currentAccount))
}

func (c *IssueBindingController) createBinding(ctx *gin.Context) {
	var request dto.CreateIssueBindingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	account, _ := ctx.Get("account")
	currentAccount := account.(*domain.Account)

	binding, err := c.bindingUC.CreateIssueAndBinding(usecases.CreateIssueBindingRequest{
		Title:           request.Title,
		ProjectId:       request.ProjectId,
		SprintId:        request.SprintId,
		AssigneeId:      request.AssigneeId,
		CreateInTracker: request.CreateInTracker,
		AccountId:       currentAccount.Id,
	})
	if apiutils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, dto.SerializeIssueBinding(binding, currentAccount))
}
