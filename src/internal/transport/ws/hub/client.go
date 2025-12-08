package hub

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn       *websocket.Conn
	ID         uint          `json:"client_id"`
	Username   string        `json:"username"`
	AvatarPath string        `json:"avatar_path"`
	Rooms      map[uint]bool `json:"-"`
}

func (c *Client) SendMessage(messageData interface{}) error {
	if c.Conn != nil {
		return c.Conn.WriteJSON(messageData)
	}
	return nil
}
