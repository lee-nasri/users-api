package apperror

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

type AppError struct {
	Message        string    `json:"message"`
	Code           ErrorCode `json:"code"`
	HTTPStatusCode int       `json:"-"`
}

func (e AppError) Error() string {
	return e.Message
}

func IsAppError(err error) (e AppError, ok bool) {
	ok = errors.As(err, &e)
	return
}

func NewErrDefault(err error) error {
	return AppError{
		Code:           ErrInternalServerError,
		Message:        err.Error(),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func NewErrInternal() error {
	return AppError{
		Code:           ErrInternalServerError,
		Message:        Message[ErrInternalServerError],
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func NewInvalidRequest() error {
	return AppError{
		Code:           ErrInvalidRequest,
		Message:        Message[ErrInvalidRequest],
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewInvalidRequestFromErr(err error) error {
	errMsg := err.Error()

	// err.Error() of echo.HTTPError contains code and message. But we only need the message.
	if httpError, ok := err.(*echo.HTTPError); ok {
		if msg, ok := httpError.Message.(string); ok {
			errMsg = msg
		}
	}

	return AppError{
		Code:           ErrInvalidRequest,
		Message:        errMsg,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewInvalidRequestWithMsg(Msg string) error {
	return AppError{
		Code:           ErrInvalidRequest,
		Message:        Msg,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewErrUserNotFound() error {
	return AppError{
		Code:           ErrUserNotFound,
		Message:        Message[ErrUserNotFound],
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewErrUserIDAlreadyExist() error {
	return AppError{
		Code:           ErrUserIDAleadyExist,
		Message:        Message[ErrUserIDAleadyExist],
		HTTPStatusCode: http.StatusConflict,
	}
}
