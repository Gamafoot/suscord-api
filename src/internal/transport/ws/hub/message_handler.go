package hub

import (
	"suscord/internal/transport/ws/hub/dto"
)

// handleClientMessage обрабатывает сообщения от клиента
func (hub *Hub) handleClientMessage(client *Client, message *dto.ClientMessage) error {
	switch message.Type {
	default:
		return client.SendMessage(&dto.ResponseMessage{
			Type: "error",
			Data: map[string]interface{}{"message": "unknown message type"},
		})
	}
}
