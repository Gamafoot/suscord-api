package event

type ChatUserLeave struct {
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
}

func NewChatUserLeave(chatID, userID uint) ChatUserLeave {
	return ChatUserLeave{
		ChatID: chatID,
		UserID: userID,
	}
}

func (e ChatUserLeave) EventName() string {
	return "chat.user_leave"
}

func (e ChatUserLeave) AggregateID() uint {
	return e.ChatID
}
