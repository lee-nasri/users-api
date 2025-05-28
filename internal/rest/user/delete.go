package user

import (
	"github.com/labstack/echo/v4"

	"users-api/domain"
	"users-api/pkg/httpserver"
	"users-api/pkg/logx"
)

func (h *Handler) DeleteUser(ctx echo.Context) error {
	c := ctx.Request().Context()
	id := ctx.Param("id")

	res, err := h.userSvc.DeleteUser(c, id)
	if err != nil {
		logx.Error(c, err, "error deleting user")
		return httpserver.NewErrorResponse(ctx, err)
	}

	return httpserver.NewSuccessResponse(ctx, domain.DeleteUserResponse{
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
