package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/infra/hub"
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
	hub       *hub.SprintHub                 `di.inject:"SprintHub"`
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

	sprintId, err := c.bindingUC.DeleteBinding(uint(id))
	if apiutils.HandleException(ctx, err) {
		return
	}

	c.hub.Broadcast(sprintId, hub.WSEvent{
		Type: hub.EventBindingDeleted,
		Data: hub.BindingDeletedEvent{ID: uint(id), SprintID: sprintId},
	})

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

	serialized := dto.SerializeIssueBinding(binding, currentAccount)
	c.hub.Broadcast(binding.Sprint.Id, hub.WSEvent{
		Type: hub.EventBindingUpdated,
		Data: hub.BindingUpdatedEvent{ID: binding.Id, AccountID: currentAccount.Id},
	})

	ctx.JSON(http.StatusOK, serialized)
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

	sprintId, err := c.bindingUC.SaveOrdering(ordering)
	if apiutils.HandleException(ctx, err) {
		return
	}

	if sprintId > 0 {
		c.hub.Broadcast(sprintId, hub.WSEvent{
			Type: hub.EventBindingOrdering,
			Data: request,
		})
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

	serialized := dto.SerializeIssueBinding(binding, currentAccount)
	c.hub.Broadcast(binding.Sprint.Id, hub.WSEvent{
		Type: hub.EventBindingCreated,
		Data: hub.BindingCreatedEvent{ID: binding.Id, AccountID: currentAccount.Id},
	})

	ctx.JSON(http.StatusCreated, serialized)
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

	binding, oldSprintId, err := c.bindingUC.MoveBinding(usecases.MoveIssueBindingRequest{
		Id:       uint(id),
		SprintId: request.SprintId,
	}, currentAccount)
	if apiutils.HandleException(ctx, err) {
		return
	}

	// уведомить старый спринт об удалении
	c.hub.Broadcast(oldSprintId, hub.WSEvent{
		Type: hub.EventBindingDeleted,
		Data: hub.BindingDeletedEvent{ID: uint(id), SprintID: oldSprintId},
	})
	// уведомить новый спринт о создании
	serialized := dto.SerializeIssueBinding(binding, currentAccount)
	c.hub.Broadcast(binding.Sprint.Id, hub.WSEvent{
		Type: hub.EventBindingCreated,
		Data: hub.BindingCreatedEvent{ID: binding.Id, AccountID: currentAccount.Id},
	})

	ctx.JSON(http.StatusOK, serialized)
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
		Title:      request.Title,
		ProjectId:  request.ProjectId,
		SprintId:   request.SprintId,
		AssigneeId: request.AssigneeId,
	})
	if apiutils.HandleException(ctx, err) {
		return
	}

	serialized := dto.SerializeIssueBinding(binding, currentAccount)
	c.hub.Broadcast(binding.Sprint.Id, hub.WSEvent{
		Type: hub.EventBindingCreated,
		Data: hub.BindingCreatedEvent{ID: binding.Id, AccountID: currentAccount.Id},
	})

	ctx.JSON(http.StatusCreated, serialized)
}
