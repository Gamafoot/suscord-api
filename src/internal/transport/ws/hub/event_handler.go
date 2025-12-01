package hub

import (
	"suscord/internal/domain/eventbus"
	"suscord/internal/transport/ws/hub/dto"
)

func (hub *Hub) RegisterEventSubscribers(bus eventbus.Bus) {
	bus.Subscribe(eventbus.EventMessageCreated, hub.onMessageCreated)
	bus.Subscribe(eventbus.EventJoinedRoom, hub.onUserJoined)
	bus.Subscribe(eventbus.EventLeftRoom, hub.onUserLeft)
}

func (hub *Hub) onMessageCreated(event eventbus.Event) {
	message := event.(eventbus.MessageCreated)
	hub.broadcast <- &dto.ResponseMessage{
		Type:   "message",
		RoomID: message.ChatID,
		Data:   message,
	}
}

func (hub *Hub) onUserJoined(event eventbus.Event) {
	joined := event.(eventbus.JoinedRoom)
	hub.broadcast <- &dto.ResponseMessage{
		Type:   "user_joined",
		RoomID: joined.RoomID,
		Data:   joined.User,
	}
}

func (hub *Hub) onUserLeft(event eventbus.Event) {
	left := event.(eventbus.LeftRoom)
	hub.broadcast <- &dto.ResponseMessage{
		Type:   "user_left",
		RoomID: left.RoomID,
		Data:   left.User,
	}
}
