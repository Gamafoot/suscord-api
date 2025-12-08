package dto

import "suscord/internal/domain/eventbus/events"

type NewUserInChat struct {
	ChatID uint  `json:"chat_id"`
	User   *User `json:"user"`
}

func (NewUserInChat) EventName() string {
	return events.EventNewUserInChat
}
