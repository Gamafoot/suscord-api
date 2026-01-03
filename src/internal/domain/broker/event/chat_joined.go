package event

type ChatUserJoined struct {
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
}

func NewChatUserJoined(chatID, userID uint) ChatUserJoined {
	return ChatUserJoined{
		ChatID: chatID,
		UserID: userID,
	}
}

func (e ChatUserJoined) EventName() string {
	return "chat.user_joined"
}

func (e ChatUserJoined) AggregateID() uint {
	return e.ChatID
}
