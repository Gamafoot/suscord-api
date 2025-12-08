package hub

import (
	"suscord/internal/domain/eventbus"
	eventDTO "suscord/internal/domain/eventbus/dto"
	"suscord/internal/domain/eventbus/events"
	"suscord/internal/transport/ws/hub/dto"
)

func (hub *Hub) RegisterEventSubscribers(bus eventbus.Bus) {
	bus.Subscribe(events.EventMessageCreate, hub.onMessageCreate)
	bus.Subscribe(events.EventMessageUpdate, hub.onMessageUpdate)
	bus.Subscribe(events.EventMessageDelete, hub.onMessageDelete)
	bus.Subscribe(events.EventInviteToChat, hub.onInviteToRoom)
	bus.Subscribe(events.EventJoinedGroupChat, hub.onJoinedGroupChat)
	bus.Subscribe(events.EventJoinedPrivateChat, hub.onJoinedPrivateChat)
	bus.Subscribe(events.EventUpdateGroupChat, hub.onUpdateGroupChat)
	bus.Subscribe(events.EventNewUserInChat, hub.onNewUserInChat)
	bus.Subscribe(events.EventLeftChat, hub.onLeftChat)
	bus.Subscribe(events.EventDeleteChat, hub.onDeleteChat)
}

func (hub *Hub) onMessageCreate(event eventbus.Event) {
	data := event.(*eventDTO.Message)
	hub.broadcastToRoom(data.ChatID, &dto.ResponseMessage{
		Type:   "message",
		ChatID: data.ChatID,
		Data:   data,
	})
}

func (hub *Hub) onMessageUpdate(event eventbus.Event) {
	data := event.(*eventDTO.Message)
	hub.broadcastToRoomExcept(data.ChatID, data.UserID, &dto.ResponseMessage{
		Type:   "message_update",
		ChatID: data.ChatID,
		Data:   data,
	})
}

func (hub *Hub) onMessageDelete(event eventbus.Event) {
	data := event.(*eventDTO.MessageDelete)
	hub.broadcastToRoomExcept(data.ChatID, data.ExceptUserID, &dto.ResponseMessage{
		Type:   "message_delete",
		ChatID: data.ChatID,
		Data:   data,
	})
}

func (hub *Hub) onInviteToRoom(event eventbus.Event) {
	data := event.(*eventDTO.InviteToRoom)
	if client, exists := hub.clients[data.UserID]; exists {
		client.SendMessage(&dto.ResponseMessage{
			Type: "invite_to_chat",
			Data: map[string]string{
				"code": data.Code,
			},
		})
	}
}

func (hub *Hub) onJoinedGroupChat(event eventbus.Event) {
	data := event.(*eventDTO.JoinedGroupChat)
	if client, exists := hub.clients[data.UserID]; exists {
		hub.joinRoom(data.Chat.ID, client)
		hub.broadcastToRoomExcept(data.Chat.ID, data.UserID, &dto.ResponseMessage{
			Type:   "joined_chat",
			ChatID: data.Chat.ID,
			Data:   data,
		})
	}
}

func (hub *Hub) onJoinedPrivateChat(event eventbus.Event) {
	data := event.(*eventDTO.JoinedPrivateChat)
	if client, exists := hub.clients[data.UserID]; exists {
		hub.joinRoom(data.Chat.ID, client)
		if !data.DontSend {
			client.SendMessage(&dto.ResponseMessage{
				Type: "joined_chat",
				Data: data,
			})
		}
	}
}

func (hub *Hub) onUpdateGroupChat(event eventbus.Event) {
	data := event.(*eventDTO.UpdateGroupChat)
	hub.broadcastToRoomExcept(data.Chat.ID, data.ExceptUserID, &dto.ResponseMessage{
		Type: "update_group_chat",
		Data: data,
	})
}

func (hub *Hub) onNewUserInChat(event eventbus.Event) {
	data := event.(*eventDTO.NewUserInChat)
	if client, exists := hub.clients[data.User.ID]; exists {
		hub.joinRoom(data.ChatID, client)
		hub.broadcastToRoomExcept(data.ChatID, data.User.ID, &dto.ResponseMessage{
			Type:   "new_user_in_chat",
			ChatID: data.ChatID,
			Data:   data,
		})
	}
}

func (hub *Hub) onLeftChat(event eventbus.Event) {
	data := event.(*eventDTO.LeftChat)
	hub.leaveRoom(data.ChatID, data.UserID)
	hub.broadcastToRoomExcept(data.ChatID, data.UserID, &dto.ResponseMessage{
		Type:   "user_left",
		ChatID: data.ChatID,
		Data:   data,
	})
}

func (hub *Hub) onDeleteChat(event eventbus.Event) {
	data := event.(*eventDTO.DeleteChat)
	hub.broadcastToRoomExcept(data.ChatID, data.ExceptUserID, &dto.ResponseMessage{
		Type: "delete_chat",
		Data: data,
	})
	hub.deleteRoom(data.ChatID)
}
