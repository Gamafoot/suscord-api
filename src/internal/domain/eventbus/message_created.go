package eventbus

import "time"

const EventMessageCreated = "MessageCreated"

type MessageCreated struct {
	ID        uint      `json:"id"`
	ChatID    uint      `json:"chat_id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MessageCreated) Name() string {
	return EventMessageCreated
}
