package dto

import "suscord/internal/domain/entity"

type Attachment struct {
	ID        uint   `json:"id"`
	MessageID uint   `json:"message_id"`
	FilePath  string `json:"file_path"`
	FileSize  int64  `json:"file_size"`
	MimeType  string `json:"mime_type"`
}

func NewAttachmentResponse(attachment *entity.Attachment) *Attachment {
	return &Attachment{
		ID:        attachment.ID,
		MessageID: attachment.MessageID,
		FilePath:  attachment.FilePath,
		FileSize:  attachment.FileSize,
		MimeType:  attachment.MimeType,
	}
}

type CreateAttachmentRequest struct {
	MessageID uint `json:"message_id" validate:"gt=0"`
}
