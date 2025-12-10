package hub

import (
	"encoding/json"
	"suscord/internal/transport/ws/hub/dto"
)

type WebRTCSignal struct {
	ChatID    uint        `json:"chatId"`
	Offer     interface{} `json:"offer,omitempty"`
	Answer    interface{} `json:"answer,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
}

// handleClientMessage обрабатывает сообщения от клиента
func (hub *Hub) handleClientMessage(client *Client, message *dto.ClientMessage) error {
	switch message.Type {
	case "call-offer", "call-answer", "ice-candidate", "call-ended", "call-declined":
		return hub.handleWebRTCSignal(client, message)
	default:
		return client.SendMessage(&dto.ResponseMessage{
			Type: "error",
			Data: map[string]interface{}{"message": "unknown message type"},
		})
	}
}

func (hub *Hub) handleWebRTCSignal(client *Client, message *dto.ClientMessage) error {
	var dataStr string
	if err := json.Unmarshal(message.Data, &dataStr); err != nil {
		return err
	}
	
	var signal WebRTCSignal
	if err := json.Unmarshal([]byte(dataStr), &signal); err != nil {
		return err
	}

	// Проверяем, что клиент является членом чата
	hub.mutex.RLock()
	isInRoom := client.Rooms[signal.ChatID]
	hub.mutex.RUnlock()

	if !isInRoom {
		return client.SendMessage(&dto.ResponseMessage{
			Type: "error",
			Data: map[string]interface{}{"message": "not a member of this chat"},
		})
	}

	// Пробрасываем сигнал всем участникам чата, кроме отправителя
	responseData := map[string]interface{}{
		"chatId": signal.ChatID,
	}

	switch message.Type {
	case "call-offer":
		responseData["offer"] = signal.Offer
	case "call-answer":
		responseData["answer"] = signal.Answer
	case "ice-candidate":
		responseData["candidate"] = signal.Candidate
	}

	hub.mutex.RLock()
	room := hub.rooms[signal.ChatID]
	hub.mutex.RUnlock()

	for clientID := range room {
		if clientID != client.ID {
			hub.mutex.RLock()
			targetClient := hub.clients[clientID]
			hub.mutex.RUnlock()

			if targetClient != nil {
				targetClient.SendMessage(&dto.ResponseMessage{
					Type: message.Type,
					Data: responseData,
				})
			}
		}
	}

	return nil
}
