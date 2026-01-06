package message

import (
	"suscord/internal/domain/broker/event"
	"suscord/internal/domain/broker/message/model"
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
	return event.OnMessageCreated
}
