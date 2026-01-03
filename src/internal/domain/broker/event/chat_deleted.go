package event

type ChatDeleted struct {
	ID uint `json:"id"`
}

func NewChatDeleted(id uint) ChatDeleted {
	return ChatDeleted{ID: id}
}

func (e ChatDeleted) EventName() string {
	return "chat.deleted"
}

func (e ChatDeleted) AggregateID() uint {
	return e.ID
}
