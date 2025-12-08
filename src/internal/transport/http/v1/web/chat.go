package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitChatRoutes(route *echo.Group) {
	route.GET("/", h.ChatPage)
	route.GET("/chats/:id", h.ChatPage)
	route.GET("/get-session", h.GetSession)
}

func (h *handler) ChatPage(c echo.Context) error {
	return c.Render(http.StatusOK, "chat.html", nil)
}
