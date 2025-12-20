package hub

import (
	"context"
	"suscord/internal/domain/entity"
	domainError "suscord/internal/domain/errors"
	"suscord/internal/transport/ws/hub/dto"
)

// Room management - управление комнатами и клиентами
func (hub *Hub) joinChatRoom(chatID uint, client *Client) error {
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
	if _, exists := hub.chatRooms[chatID]; !exists {
		hub.chatRooms[chatID] = make(map[uint]bool)
	}

	// Добавляем клиента в комнату
	hub.chatRooms[chatID][client.ID] = true
	client.Rooms[chatID] = true

	return nil
}

func (hub *Hub) leaveChatRoom(roomID, clientID uint) {
	hub.mutex.Lock()

	// Удаляем клиента из комнаты
	if room, exists := hub.chatRooms[roomID]; exists {
		delete(room, clientID)
		if len(room) == 0 {
			delete(hub.chatRooms, roomID)
		}
	}

	// Удаляем комнату у клиента
	if client, exists := hub.clients[clientID]; exists {
		delete(client.Rooms, roomID)
	}

	hub.mutex.Unlock()

	// Уведомляем других участников
	hub.broadcastToChatRoomExcept(roomID, clientID, &dto.ResponseMessage{
		Type: "user_left",
		Data: map[string]interface{}{"user_id": clientID},
	})
}

func (hub *Hub) deleteChatRoom(roomID uint) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.chatRooms[roomID] = make(map[uint]bool)

	for clientID := range hub.clients {
		delete(hub.clients[clientID].Rooms, roomID)
	}
}

func (hub *Hub) joinToUserChatRooms(client *Client, chats []*entity.Chat) error {
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
		if _, exists := hub.chatRooms[chat.ID]; !exists {
			hub.chatRooms[chat.ID] = make(map[uint]bool)
		}

		// Добавляем клиента в комнату
		hub.chatRooms[chat.ID][client.ID] = true
		client.Rooms[chat.ID] = true
	}

	return nil
}

func (hub *Hub) broadcastToChatRoom(roomID uint, message interface{}) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.chatRooms[roomID]; exists {
		for userID := range room {
			if client, exists := hub.clients[userID]; exists {
				client.SendMessage(message)
			}
		}
	}
}

func (hub *Hub) broadcastToChatRoomExcept(roomID, exceptUserID uint, message any) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.chatRooms[roomID]; exists {
		for userID := range room {
			if userID != exceptUserID {
				if client, exists := hub.clients[userID]; exists {
					client.SendMessage(message)
				}
			}
		}
	}
}

func (hub *Hub) broadcastToSFURoom(roomID uint, message any) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.sfuRooms[roomID]; exists {
		for userID := range room {
			if client, exists := hub.clients[userID]; exists {
				client.SendMessage(message)
			}
		}
	}
}

func (hub *Hub) broadcastToSFURoomExcept(roomID, exceptUserID uint, message any) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	if room, exists := hub.sfuRooms[roomID]; exists {
		for userID := range room {
			if userID != exceptUserID {
				if client, exists := hub.clients[userID]; exists {
					client.SendMessage(message)
				}
			}
		}
	}
}
