package event

type ChatInvite struct {
	ChatID uint `json:"chat_id"`
	UserID uint `json:"user_id"`
}

func NewChatInvite(chatID, userID uint) ChatInvite {
	return ChatInvite{
		ChatID: chatID,
		UserID: userID,
	}
}

func (e ChatInvite) EventName() string {
	return "chat.invite"
}

func (e ChatInvite) AggregateID() uint {
	return e.ChatID
}
