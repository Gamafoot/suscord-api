package dto

import "suscord/internal/domain/eventbus/events"

type UpdateGroupChat struct {
	Chat         *Chat `json:"chat"`
	ExceptUserID uint  `json:"-"`
}

func (UpdateGroupChat) EventName() string {
	return events.EventUpdateGroupChat
}
