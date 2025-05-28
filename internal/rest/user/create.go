package user

import (
	"github.com/labstack/echo/v4"

	"users-api/domain"
	"users-api/pkg/apperror"
	"users-api/pkg/httpserver"
	"users-api/pkg/logx"
)

func (h *Handler) CreateUser(ctx echo.Context) error {
	c := ctx.Request().Context()
	user := new(domain.CreateUserRequest)

	if err := ctx.Bind(user); err != nil {
		logx.Error(c, err, "error binding user data")
		return httpserver.NewErrorResponse(ctx, apperror.NewInvalidRequest())
	}

	if err := h.validateBodyParser(user); err != nil {
		logx.Error(c, err, "error validating user data")
		return httpserver.NewErrorResponse(ctx, apperror.NewInvalidRequest())
	}

	res, err := h.userSvc.CreateUser(c, *user)
	if err != nil {
		logx.Error(c, err, "error creating user")
		return httpserver.NewErrorResponse(ctx, err)
	}

	return httpserver.NewSuccessResponse(ctx, domain.CreateUsersResponse{
		Data: &domain.UserResponse{
			ID:        res.ID,
			SurName:   res.Surname,
			LastName:  res.Lastname,
			Age:       res.Age,
			Email:     res.Email,
			Phone:     res.Phone,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,

			// New Fields
			FatherName: res.FatherName,
			MotherName: res.MotherName,
		},
	})

}
