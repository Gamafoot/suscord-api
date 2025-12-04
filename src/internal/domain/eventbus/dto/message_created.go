package dto

import (
	"suscord/internal/domain/eventbus/events"
	"time"
)

type MessageCreated struct {
	ID          uint          `json:"id"`
	ChatID      uint          `json:"chat_id"`
	UserID      uint          `json:"user_id"`
	Type        string        `json:"type"`
	Content     string        `json:"content"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Attachments []*Attachment `json:"attachments"`
}

func (MessageCreated) Name() string {
	return events.EventMessageCreated
}
