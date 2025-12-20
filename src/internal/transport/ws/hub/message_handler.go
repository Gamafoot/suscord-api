package hub

import (
	"suscord/internal/transport/ws/hub/dto"
	"suscord/internal/transport/ws/hub/mapper"

	pkgErrors "github.com/pkg/errors"
)

func (hub *Hub) handleClientMessage(client *Client, message *dto.ClientMessage) error {
	switch message.Type {
	case "call-invite":
		if err := hub.joinSFURoom(message.ChatID, client); err != nil {
			return err
		}
		hub.broadcastToChatRoomExcept(message.ChatID, client.ID, message)

	case "call-accept":
		if err := hub.joinSFURoom(message.ChatID, client); err != nil {
			return err
		}
		hub.broadcastToSFURoomExcept(message.ChatID, client.ID, &dto.ResponseMessage{
			Type: "call-accept",
			Data: mapper.NewClient(&client.Client, hub.cfg.Media.Url),
		})

		clients, err := hub.clientsSFURoom(message.ChatID)
		if err != nil {
			return err
		}
		client.SendMessage(&dto.ResponseMessage{
			Type: "call-clients",
			Data: map[string]any{
				"clients": clients,
			},
		})

	case "call-reject":
		hub.broadcastToSFURoomExcept(message.ChatID, client.ID, message)

	case "call-leave":
		ok := hub.isMemberOfSFURoom(message.ChatID, client.ID)
		if !ok {
			return pkgErrors.New("you are not member of room")
		}

		if err := hub.leaveSFURoom(message.ChatID, client); err != nil {
			return err
		}

		clients, err := hub.clientsSFURoom(message.ChatID)
		if err != nil {
			return err
		}

		hub.broadcastToSFURoomExcept(message.ChatID, client.ID, &dto.ResponseMessage{
			Type: "call-leave",
			Data: map[string]any{
				"clients": clients,
				"user_id": client.ID,
			},
		})

	case "call-ended":
		hub.broadcastToSFURoom(message.ChatID, message)

	case "call-stream":
		hub.broadcastToSFURoomExcept(message.ChatID, client.ID, message)

	default:
		return client.SendMessage(&dto.ResponseMessage{
			Type: "error",
			Data: map[string]interface{}{"message": "unknown message type", "data": message},
		})
	}

	return nil
}
