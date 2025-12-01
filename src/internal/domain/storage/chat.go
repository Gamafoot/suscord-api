package storage

import (
	"context"
	"suscord/internal/domain/entity"
)

type ChatStorage interface {
	GetByID(ctx context.Context, chatID uint) (*entity.Chat, error)
	GetUserChat(ctx context.Context, chatID, userID uint) (*entity.Chat, error)
	GetUserChats(ctx context.Context, userID uint) ([]*entity.Chat, error)
	GetPrivateChatID(ctx context.Context, userID, friendID uint) (uint, error)
	Create(ctx context.Context, data *entity.CreateChat) (uint, error)
	Update(ctx context.Context, chatID uint, data *entity.UpdateChat) error
	Delete(ctx context.Context, chatID uint) error
}
