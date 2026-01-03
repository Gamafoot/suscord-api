package event

type MessageDeleted struct {
	ChatID    uint
	MessageID uint
}

func NewMessageDeleted(chatID, messageID uint) MessageDeleted {
	return MessageDeleted{
		ChatID:    chatID,
		MessageID: messageID,
	}
}

func (e MessageDeleted) EventName() string {
	return "chat.message.deleted"
}

func (e MessageDeleted) AggregateID() uint {
	return e.ChatID
}
