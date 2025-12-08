package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewMessage(event string, message *entity.Message) *dto.Message {
	attachments := make([]*dto.Attachment, len(message.Attachments))
	for i, attachment := range message.Attachments {
		attachments[i] = &dto.Attachment{
			ID:       attachment.ID,
			FilePath: attachment.FilePath,
			FileSize: attachment.FileSize,
			MimeType: attachment.MimeType,
		}
	}

	return &dto.Message{
		Event:       event,
		ID:          message.ID,
		ChatID:      message.ChatID,
		UserID:      message.UserID,
		Content:     message.Content,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
		Attachments: attachments,
	}
}
