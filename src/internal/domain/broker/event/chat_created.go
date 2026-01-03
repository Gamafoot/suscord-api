package event

import (
	"suscord/internal/domain/broker/event/model"
	"suscord/internal/domain/entity"
)

type ChatCreated struct {
	model.Chat
}

func NewChatCreated(chat *entity.Chat, mediaURL string) ChatCreated {
	return ChatCreated{
		Chat: model.NewChat(chat, mediaURL),
	}
}

func (e ChatCreated) EventName() string {
	return "chat.created"
}

func (e ChatCreated) AggregateID() uint {
	return e.ID
}
