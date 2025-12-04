package model

import "time"

type Message struct {
	ID        uint
	ChatID    uint
	UserID    uint
	Type      string `gorm:"type:varchar(50)"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Chat        Chat
	User        User
	Attachments []*Attachment `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Message) TableName() string {
	return "messages"
}
