package model

import "time"

type Message struct {
	ID        uint
	ChatID    uint
	UserID    uint
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Chat Chat
	User User
}

func (m *Message) TableName() string {
	return "messages"
}
