package message

import (
	"suscord/internal/domain/broker/event"
	"suscord/internal/domain/broker/message/model"
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
	return event.OnChatUpdated
}
