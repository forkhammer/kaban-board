package utils

import (
	"errors"
	"main/internal/domain"
	"main/internal/interfaces/api/dto"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func HandleException(ctx *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if validErr, ok := errors.AsType[*domain.ValidationError](err); ok {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: validErr.Error()})
		return true
	}

	if foundErr, ok := errors.AsType[*domain.NotFoundError](err); ok {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: foundErr.Error()})
		return true
	}

	sentry.CaptureException(err)
	ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	return true
}
