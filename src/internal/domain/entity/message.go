package entity

import (
	"time"
)

type Message struct {
	ID          uint
	ChatID      uint
	UserID      uint
	Type        string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Attachments []*Attachment
}

type GetMessagesInput struct {
	ChatID        uint
	UserID        uint
	LastMessageID uint
	Limit         int
}

type CreateMessage struct {
	Type    string
	Content string
}

type UpdateMessage struct {
	Content string
}
