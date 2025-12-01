package model

type MessageAttachment struct {
	ID        uint
	MessageID uint
	FileURL   string `gorm:"varchar(255)"`
	FileSize  int
	MimeType  string `gorm:"varchar(100)"`

	Message Message
}

func (m *MessageAttachment) TableName() string {
	return "message_attachments"
}
