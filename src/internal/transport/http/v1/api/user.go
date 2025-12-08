package api

import (
	"errors"
	"net/http"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/dto"
	"suscord/internal/transport/mapper"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitUserRoutes(route *echo.Group) {
	route.GET("/users/:user_id", h.GetUserInfo)
	route.GET("/users/me", h.AboutMe)
	route.GET("/users", h.SearchUsers)
}

func (h *handler) GetUserInfo(c echo.Context) error {
	userID, err := utils.GetUIntFromParam(c, "user_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, "user_id must be digit")
	}

	user, err := h.service.User().GetByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	userDTO := dto.User{
		ID:         user.ID,
		Username:   user.Username,
		AvatarPath: user.AvatarPath,
	}

	return c.JSON(http.StatusOK, userDTO)
}

func (h *handler) AboutMe(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	user, err := h.service.User().GetByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	userDTO := dto.Me{
		ID:         user.ID,
		Username:   user.Username,
		AvatarPath: user.AvatarPath,
		FriendCode: user.FriendCode,
	}

	return c.JSON(http.StatusOK, userDTO)
}

func (h *handler) SearchUsers(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	searchPattern := c.QueryParam("search")

	users, err := h.service.User().SearchUsers(c.Request().Context(), userID, searchPattern)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	result := make([]*dto.User, len(users))

	for i, user := range users {
		result[i] = mapper.NewUser(user)
	}

	return c.JSON(http.StatusOK, result)
}
