package event

import (
	"suscord/internal/domain/broker/event/model"
	"suscord/internal/domain/entity"
)

type ChatUpdated struct {
	model.Chat
}

func NewChatUpdated(chat *entity.Chat, mediaURL string) ChatUpdated {
	return ChatUpdated{
		Chat: model.NewChat(chat, mediaURL),
	}
}

func (e ChatUpdated) EventName() string {
	return "chat.udpated"
}

func (e ChatUpdated) AggregateID() uint {
	return e.ID
}
