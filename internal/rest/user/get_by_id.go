package user

import (
	"github.com/labstack/echo/v4"

	"users-api/domain"
	"users-api/pkg/httpserver"
	"users-api/pkg/logx"
)

func (h *Handler) GetUserByID(ctx echo.Context) error {
	c := ctx.Request().Context()
	userID := ctx.Param("id")

	res, err := h.userSvc.GetUserByUserID(c, userID)
	if err != nil {
		logx.Error(c, err, "error getting user by id")
		return httpserver.NewErrorResponse(ctx, err)
	}

	return httpserver.NewSuccessResponse(ctx, domain.GetUserResponse{
		Data: &domain.UserResponse{
			ID:        res.ID,
			SurName:   res.Surname,
			LastName:  res.Lastname,
			Age:       res.Age,
			Email:     res.Email,
			Phone:     res.Phone,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}
