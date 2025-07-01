package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReleaseController struct {
	releaseUC *usecases.ReleaseUseCases `di.inject:"ReleaseUseCases"`
}

func (c *ReleaseController) RegisterRoutes(router *gin.Engine) error {
	releaseRoutes := router.Group("/")
	releaseRoutes.GET("/release", c.getReleases)
	releaseRoutes.GET("/release/:id", c.getRelease)
	return nil
}

func (c *ReleaseController) getReleases(ctx *gin.Context) {
	var request dto.GetReleasesRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	releases, err := c.releaseUC.GetReleases(&queries.ReleaseFilter{
		Search:    request.Search,
		ProjectId: request.Project,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeReleases(releases))
}

func (c *ReleaseController) getRelease(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	release, err := c.releaseUC.GetRelease(uint(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeRelease(release))
}
