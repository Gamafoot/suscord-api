package eventbus

const EventLeftRoom = "LeftRoom"

type LeftRoom struct {
	RoomID uint `json:"room_id"`
	User   User `json:"user"`
}

func (LeftRoom) Name() string {
	return EventLeftRoom
}
