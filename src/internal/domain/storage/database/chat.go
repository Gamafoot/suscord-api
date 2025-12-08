package database

import (
	"context"
	"suscord/internal/domain/entity"
)

type ChatStorage interface {
	GetByID(ctx context.Context, chatID uint) (*entity.Chat, error)
	GetUserChat(ctx context.Context, chatID, userID uint) (*entity.Chat, error)
	GetUserChats(ctx context.Context, userID uint) ([]*entity.Chat, error)
	SearchUserChats(ctx context.Context, userID uint, searchPattern string) ([]*entity.Chat, error)
	Create(ctx context.Context, data *entity.CreateChat) (uint, error)
	Update(ctx context.Context, chatID uint, data *entity.UpdateChat) error
	Delete(ctx context.Context, chatID uint) error
}
