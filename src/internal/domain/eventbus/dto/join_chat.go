package dto

import "suscord/internal/domain/eventbus/events"

type JoinedGroupChat struct {
	Chat     *Chat `json:"chat"`
	UserID   uint  `json:"-"`
	DontSend bool  `json:"-"`
}

func (JoinedGroupChat) EventName() string {
	return events.EventJoinedGroupChat
}
