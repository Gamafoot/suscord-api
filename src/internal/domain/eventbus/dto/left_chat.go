package dto

import "suscord/internal/domain/eventbus/events"

type LeftChat struct {
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
}

func (LeftChat) EventName() string {
	return events.EventLeftChat
}
