package hub

import (
	"context"
	"suscord/internal/domain/entity"
	domainError "suscord/internal/domain/errors"
	"suscord/internal/transport/ws/hub/dto"
)

// Room management - управление комнатами и клиентами
func (hub *Hub) joinRoom(roomID uint, client *Client) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), hub.cfg.Server.Timeout)
	defer cancel()

	// Проверяем права доступа
	ok, err := hub.storage.ChatMember().IsMemberOfChat(ctx, roomID, client.ID)
	if err != nil {
		return err
	}

	if !ok {
		return client.SendMessage(&dto.ResponseMessage{
			Type: "join_room_error",
			Data: map[string]interface{}{"message": "You are not member of this room"},
		})
	}

	// Создаем комнату если не существует
	if _, exists := hub.rooms[roomID]; !exists {
		hub.rooms[roomID] = make(map[uint]bool)
	}

	// Добавляем клиента в комнату
	hub.rooms[roomID][client.ID] = true
	client.Rooms[roomID] = true

	// Отправляем подтверждение
	return client.SendMessage(&dto.ResponseMessage{
		Type: "join_room_ok",
		Data: map[string]interface{}{"room_id": roomID},
	})
}

func (hub *Hub) leaveRoom(roomID, clientID uint) {
	// Удаляем клиента из комнаты
	if room, exists := hub.rooms[roomID]; exists {
		delete(room, clientID)
		if len(room) == 0 {
			delete(hub.rooms, roomID)
		}
	}

	// Удаляем комнату у клиента
	if client, exists := hub.clients[clientID]; exists {
		delete(client.Rooms, roomID)
	}

	// Уведомляем других участников
	go hub.broadcastToRoomExcept(roomID, clientID, &dto.ResponseMessage{
		Type: "user_left",
		Data: map[string]interface{}{"user_id": clientID},
	})
}

func (hub *Hub) joinToUserRooms(client *Client, chats []*entity.Chat) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), hub.cfg.Server.Timeout)
	defer cancel()

	for _, chat := range chats {
		ok, err := hub.storage.ChatMember().IsMemberOfChat(ctx, client.ID, chat.ID)
		if err != nil {
			return err
		}

		if !ok {
			return domainError.ErrUserIsNotMemberOfChat
		}

		// Создаем комнату если не существует
		if _, exists := hub.rooms[chat.ID]; !exists {
			hub.rooms[chat.ID] = make(map[uint]bool)
		}

		// Добавляем клиента в комнату
		hub.rooms[chat.ID][client.ID] = true
		client.Rooms[chat.ID] = true
	}

	return nil
}

// Broadcasting methods
func (hub *Hub) broadcastToRoom(roomID uint, message interface{}) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.rooms[roomID]; exists {
		for userID := range room {
			if client, exists := hub.clients[userID]; exists {
				client.SendMessage(message)
			}
		}
	}
}

func (hub *Hub) broadcastToRoomExcept(roomID, exceptUserID uint, message *dto.ResponseMessage) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.rooms[roomID]; exists {
		for userID := range room {
			if userID != exceptUserID {
				if client, exists := hub.clients[userID]; exists {
					go client.SendMessage(message)
				}
			}
		}
	}
}
