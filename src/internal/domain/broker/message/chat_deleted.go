package message

import "suscord/internal/domain/broker/event"

type ChatDeleted struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewChatDeleted(id uint, name string) ChatDeleted {
	return ChatDeleted{ID: id, Name: name}
}

func (e ChatDeleted) EventName() string {
	return event.OnChatDeleted
}
