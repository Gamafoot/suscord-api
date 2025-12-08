package database

import (
	"context"
	"suscord/internal/domain/entity"
)

type MessageStorage interface {
	GetByID(ctx context.Context, messageID uint) (*entity.Message, error)
	GetMessages(ctx context.Context, chatID, lastMessageID uint, limit int) ([]*entity.Message, error)
	GetByAttachmentID(ctx context.Context, attachmentID uint) (*entity.Message, error)
	Create(ctx context.Context, userID, chatID uint, data *entity.CreateMessage) (uint, error)
	Update(ctx context.Context, messageID uint, data *entity.UpdateMessage) error
	Delete(ctx context.Context, messageID uint) error
	IsOwner(ctx context.Context, userID, messageID uint) (bool, error)
}
