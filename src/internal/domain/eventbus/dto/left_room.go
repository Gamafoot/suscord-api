package dto

import "suscord/internal/domain/eventbus/events"

type LeftRoom struct {
	RoomID uint `json:"room_id"`
	User   User `json:"user"`
}

func (LeftRoom) Name() string {
	return events.EventLeftRoom
}
