package eventbus

const EventJoinedRoom = "JoinedRoom"

type JoinedRoom struct {
	RoomID uint `json:"room_id"`
	User   User `json:"user"`
}

func (JoinedRoom) Name() string {
	return EventJoinedRoom
}
