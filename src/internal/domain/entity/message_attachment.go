package entity

type MessageAttachment struct {
	ID        uint
	MessageID uint
	FileURL   string
	FileSize  int
	MimeType  string
}
