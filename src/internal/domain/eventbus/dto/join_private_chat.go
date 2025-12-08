package dto

import "suscord/internal/domain/eventbus/events"

type JoinedPrivateChat struct {
	Chat     *Chat `json:"chat"`
	UserID   uint  `json:"-"`
	DontSend bool  `json:"-"`
}

func (JoinedPrivateChat) EventName() string {
	return events.EventJoinedPrivateChat
}
