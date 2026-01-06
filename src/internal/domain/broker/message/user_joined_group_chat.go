package message

import (
	"suscord/internal/domain/broker/event"
	"suscord/internal/domain/broker/message/model"
	"suscord/internal/domain/entity"
)

type UserJoinedGroupChat struct {
	ChatID uint        `json:"chat_id"`
	User   *model.User `json:"user"`
}

func NewUserJoinedGroupChat(chatID uint, user *entity.User, mediaURL string) UserJoinedGroupChat {
	return UserJoinedGroupChat{
		ChatID: chatID,
		User:   model.NewUser(user, mediaURL),
	}
}

func (e UserJoinedGroupChat) EventName() string {
	return event.OnUserJoinedGroupChat
}
