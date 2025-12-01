package service

import (
	"context"
	"suscord/internal/domain/entity"
)

type ChatService interface {
	GetUserChats(ctx context.Context, userID uint) ([]*entity.Chat, error)
	GetOrCreatePrivateChat(ctx context.Context, data *entity.CreatePrivateChat) (*entity.Chat, error)
	CreateGroupChat(ctx context.Context, data *entity.CreateGroupChat) (*entity.Chat, error)
	UpdateGroupChat(ctx context.Context, userID, chatID uint, data *entity.UpdateChat) (*entity.Chat, error)
	DeletePrivateChat(ctx context.Context, userID, chatID uint) error
}
