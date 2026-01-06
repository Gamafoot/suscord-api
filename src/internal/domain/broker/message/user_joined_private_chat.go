package message

import (
	"suscord/internal/domain/broker/event"
	"suscord/internal/domain/broker/message/model"
	"suscord/internal/domain/entity"
)

type UserJoinedPrivateChat struct {
	ChatID uint       `json:"chat_id"`
	User   *model.User `json:"user"`
}

func NewUserJoinedPrivateChat(chatID uint, user *entity.User, mediaURL string) UserJoinedPrivateChat {
	return UserJoinedPrivateChat{
		ChatID: chatID,
		User: model.NewUser(user, mediaURL),
	}
}

func (e UserJoinedPrivateChat) EventName() string {
	return event.OnUserJoinedPrivateChat
}
