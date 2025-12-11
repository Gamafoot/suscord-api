package dto

import (
	"path/filepath"
	"suscord/internal/domain/entity"
)

type Attachment struct {
	ID        uint   `json:"id"`
	MessageID uint   `json:"message_id"`
	FileUrl   string `json:"file_url"`
	FileSize  int64  `json:"file_size"`
	MimeType  string `json:"mime_type"`
}

func NewAttachmentResponse(attachment *entity.Attachment, mediaURL string) *Attachment {
	return &Attachment{
		ID:        attachment.ID,
		MessageID: attachment.MessageID,
		FileUrl:   filepath.Join(mediaURL, attachment.FilePath),
		FileSize:  attachment.FileSize,
		MimeType:  attachment.MimeType,
	}
}

type CreateAttachmentRequest struct {
	MessageID uint `json:"message_id" validate:"gt=0"`
}
