package service

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/storage"
)

type chatMemberService struct {
	storage storage.Storage
}

func NewChatMemberService(storage storage.Storage) *chatMemberService {
	return &chatMemberService{
		storage: storage,
	}
}

func (s *chatMemberService) GetChatMembers(ctx context.Context, chatID, memberID uint) ([]*entity.User, error) {
	ok, err := s.storage.ChatMember().IsMemberOfChat(ctx, memberID, chatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	users, err := s.storage.ChatMember().GetChatMembers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *chatMemberService) AddUserToChat(ctx context.Context, chatID, memberID, userID uint) error {
	ok, err := s.storage.ChatMember().IsMemberOfChat(ctx, memberID, chatID)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return s.storage.ChatMember().AddUserToChat(ctx, memberID, chatID)
}
