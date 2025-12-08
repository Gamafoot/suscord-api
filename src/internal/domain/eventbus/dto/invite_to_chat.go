package dto

import "suscord/internal/domain/eventbus/events"

type InviteToRoom struct {
	Code   string `json:"code"`
	UserID uint   `json:"user_id"`
	ChatID uint   `json:"chat_id"`
}

func (InviteToRoom) EventName() string {
	return events.EventInviteToChat
}
