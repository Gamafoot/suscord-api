package event

import (
	"suscord/internal/domain/broker/event/model"
	"suscord/internal/domain/entity"
)

type MessageCreated struct {
	model.Message
}

func NewMessageCreated(message *entity.Message, mediaURL string) MessageCreated {
	return MessageCreated{
		Message: model.NewMessage(message, mediaURL),
	}
}

func (e MessageCreated) EventName() string {
	return "chat.message.created"
}

func (e MessageCreated) AggregateID() uint {
	return e.ChatID
}
