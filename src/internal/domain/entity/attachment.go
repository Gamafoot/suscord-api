package entity

type Attachment struct {
	ID        uint
	MessageID uint
	FilePath  string
	FileSize  int64
	MimeType  string
}

type CreateAttachment struct {
	FilePath string
	FileSize int64
	MimeType string
}
