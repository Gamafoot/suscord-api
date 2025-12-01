package storage

import (
	"suscord/internal/domain/entity"
)

type MessageAttachmentStorage interface {
	GetByID(attachmentID int) (*entity.MessageAttachment, error)
	GetByMessageID(messageID int) (*entity.MessageAttachment, error)
	Create(attachment *entity.MessageAttachment) error
	Update(attachmentID uint, data map[string]interface{}) error
	Delete(attachmentID uint) error
}
