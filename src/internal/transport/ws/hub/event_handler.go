package hub

import (
	"suscord/internal/domain/eventbus"
	eventDTO "suscord/internal/domain/eventbus/dto"
	"suscord/internal/domain/eventbus/events"
	"suscord/internal/transport/ws/hub/dto"
)

func (hub *Hub) RegisterEventSubscribers(bus eventbus.Bus) {
	bus.Subscribe(events.EventMessageCreated, hub.onMessageCreated)
}

func (hub *Hub) onMessageCreated(event eventbus.Event) {
	message := event.(*eventDTO.MessageCreated)
	hub.broadcast <- &dto.ResponseMessage{
		Type:   "message",
		RoomID: message.ChatID,
		Data:   message,
	}
}
