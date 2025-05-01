package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"users-api/pkg/apperror"
)

func NewSuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, data)
}

func NewErrorResponse(ctx echo.Context, err error) error {
	// Response struct app error
	e, ok := apperror.IsAppError(err)
	if ok {
		return ctx.JSON(e.HTTPStatusCode, err)
	}

	// Response != app error
	return ctx.JSON(http.StatusBadRequest, apperror.NewErrDefault(err))
}
