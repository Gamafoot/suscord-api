package hub

import (
	"suscord/internal/transport/ws/hub/model"

	"github.com/gorilla/websocket"
)

type Client struct {
	model.Client
	Conn  *websocket.Conn `json:"-"`
	Rooms map[uint]bool   `json:"-"`
}

func (c *Client) SendMessage(messageData interface{}) error {
	if c.Conn != nil {
		return c.Conn.WriteJSON(messageData)
	}
	return nil
}
