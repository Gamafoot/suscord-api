package service

import (
	"context"
	"mime/multipart"
	"suscord/internal/domain/entity"
)

type MessageService interface {
	GetChatMessages(ctx context.Context, input *entity.GetMessagesInput) ([]*entity.Message, error)
	Create(ctx context.Context, userID, chatID uint, data *entity.CreateMessage, files []*multipart.FileHeader) (*entity.Message, error)
	Update(ctx context.Context, userID, messageID uint, data *entity.UpdateMessage) (*entity.Message, error)
	Delete(ctx context.Context, userID, messageID uint) error
}
