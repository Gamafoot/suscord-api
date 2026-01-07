package service

import (
	"context"
	"suscord/internal/domain/entity"
)

type ChatMemberService interface {
	GetNonMembers(ctx context.Context, chatID, memberID uint) ([]*entity.User, error)
	GetNotChatMembers(ctx context.Context, chatID, memberID uint) ([]*entity.User, error)
	IsMemberOfChat(ctx context.Context, userID, chatID uint) (bool, error)
	SendInvite(ctx context.Context, ownerID, chatID, userID uint) error
	AcceptInvite(ctx context.Context, userID uint, code string) error
	LeaveFromChat(ctx context.Context, chatID, userID uint) error
}
