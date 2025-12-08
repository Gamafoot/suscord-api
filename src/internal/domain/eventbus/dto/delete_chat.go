package dto

import "suscord/internal/domain/eventbus/events"

type DeleteChat struct {
	ChatID       uint `json:"chat_id"`
	ExceptUserID uint `json:"-"`
}

func (DeleteChat) EventName() string {
	return events.EventDeleteChat
}
