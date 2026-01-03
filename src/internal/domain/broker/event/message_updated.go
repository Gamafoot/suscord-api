package event

import (
	"suscord/internal/domain/broker/event/model"
	"suscord/internal/domain/entity"
)

type MessageUpdated struct {
	model.Message
}

func NewMessageUpdated(message *entity.Message, mediaURL string) MessageUpdated {
	return MessageUpdated{
		Message: model.NewMessage(message, mediaURL),
	}
}

func (e MessageUpdated) EventName() string {
	return "chat.message.updated"
}

func (e MessageUpdated) AggregateID() uint {
	return e.ChatID
}
