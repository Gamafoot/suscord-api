package service

import (
	"context"
	"suscord/internal/domain/entity"
)

type ChatMemberService interface {
	GetChatMembers(ctx context.Context, chatID, memberID uint) ([]*entity.User, error)
	AddUserToChat(ctx context.Context, chatID, memberID, userID uint) error
}
