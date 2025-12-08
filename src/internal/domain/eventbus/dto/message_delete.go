package dto

import (
	"suscord/internal/domain/eventbus/events"
)

type MessageDelete struct {
	ExceptUserID uint `json:"-"`
	ChatID       uint `json:"chat_id"`
	MessageID    uint `json:"message_id"`
}

func (MessageDelete) EventName() string {
	return events.EventMessageDelete
}
