package hub

import (
	"context"
	"suscord/internal/domain/entity"
	domainError "suscord/internal/domain/errors"
	"suscord/internal/transport/ws/hub/dto"
)

// Room management - управление комнатами и клиентами
func (hub *Hub) joinRoom(chatID uint, client *Client) error {
	if client == nil {
		return nil
	}

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), hub.cfg.Server.Timeout)
	defer cancel()

	// Проверяем права доступа
	ok, err := hub.storage.Database().ChatMember().IsMemberOfChat(ctx, client.ID, chatID)
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
	if _, exists := hub.rooms[chatID]; !exists {
		hub.rooms[chatID] = make(map[uint]bool)
	}

	// Добавляем клиента в комнату
	hub.rooms[chatID][client.ID] = true
	client.Rooms[chatID] = true

	return nil
}

func (hub *Hub) leaveRoom(roomID, clientID uint) {
	hub.mutex.Lock()

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

	hub.mutex.Unlock()

	// Уведомляем других участников
	hub.broadcastToRoomExcept(roomID, clientID, &dto.ResponseMessage{
		Type: "user_left",
		Data: map[string]interface{}{"user_id": clientID},
	})
}

func (hub *Hub) deleteRoom(roomID uint) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.rooms[roomID] = make(map[uint]bool)

	for clientID := range hub.clients {
		delete(hub.clients[clientID].Rooms, roomID)
	}
}

func (hub *Hub) joinToUserRooms(client *Client, chats []*entity.Chat) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), hub.cfg.Server.Timeout)
	defer cancel()

	for _, chat := range chats {
		ok, err := hub.storage.Database().ChatMember().IsMemberOfChat(ctx, client.ID, chat.ID)
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
func (hub *Hub) broadcastToRoom(chatID uint, message interface{}) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.rooms[chatID]; exists {
		for userID := range room {
			if client, exists := hub.clients[userID]; exists {
				client.SendMessage(message)
			}
		}
	}
}

func (hub *Hub) broadcastToRoomExcept(chatID, exceptUserID uint, message *dto.ResponseMessage) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.rooms[chatID]; exists {
		for userID := range room {
			if userID != exceptUserID {
				if client, exists := hub.clients[userID]; exists {
					client.SendMessage(message)
				}
			}
		}
	}
}
