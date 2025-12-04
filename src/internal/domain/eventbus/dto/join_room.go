package dto

import "suscord/internal/domain/eventbus/events"

type JoinedRoom struct {
	RoomID uint `json:"room_id"`
	User   User `json:"user"`
}

func (JoinedRoom) Name() string {
	return events.EventJoinedRoom
}
