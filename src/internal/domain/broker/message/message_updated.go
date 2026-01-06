package message

import (
	"suscord/internal/domain/broker/event"
	"suscord/internal/domain/broker/message/model"
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
	return event.OnMessageUpdated
}
