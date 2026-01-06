package message

import "suscord/internal/domain/broker/event"

type MessageDeleted struct {
	ChatID       uint
	MessageID    uint
	ExceptUserID uint
}

func NewMessageDeleted(chatID, messageID, exceptUserID uint) MessageDeleted {
	return MessageDeleted{
		ChatID:       chatID,
		MessageID:    messageID,
		ExceptUserID: exceptUserID,
	}
}

func (e MessageDeleted) EventName() string {
	return event.OnMessageDeleted
}
