package entity

import (
	"time"
)

type Message struct {
	ID        uint
	ChatID    uint
	UserID    uint
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetMessagesInput struct {
	ChatID        uint
	UserID        uint
	LastMessageID uint
	Limit         int
}

type CreateMessage struct {
	Content string
}

type UpdateMessage struct {
	Content string
}
