package message

import "suscord/internal/domain/broker/event"

type UserLeft struct {
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
}

func NewUserLeft(chatID, userID uint) UserLeft {
	return UserLeft{
		ChatID: chatID,
		UserID: userID,
	}
}

func (e UserLeft) EventName() string {
	return event.OnUserLeft
}
