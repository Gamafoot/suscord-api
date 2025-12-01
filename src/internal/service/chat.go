package service

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/storage"

	"github.com/pkg/errors"
)

type chatService struct {
	storage     storage.Storage
	fileStorage storage.FileStorage
}

func NewChatService(storage storage.Storage) *chatService {
	return &chatService{
		storage: storage,
	}
}

func (s *chatService) GetUserChats(ctx context.Context, userID uint) ([]*entity.Chat, error) {
	return s.storage.Chat().GetUserChats(ctx, userID)
}

func (s *chatService) GetOrCreatePrivateChat(ctx context.Context, input *entity.CreatePrivateChat) (*entity.Chat, error) {
	chatID, err := s.storage.Chat().GetPrivateChatID(ctx, input.UserID, input.FriendID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			if chatID == 0 {
				chatID, err = s.createPrivateChat(ctx, input)
				if err != nil {
					return nil, err
				}
			}
		}
		return nil, err
	}

	chat, err := s.storage.Chat().GetUserChat(ctx, chatID, input.UserID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) createPrivateChat(ctx context.Context, input *entity.CreatePrivateChat) (uint, error) {
	chatID, err := s.storage.Chat().Create(ctx, &entity.CreateChat{Type: "private"})
	if err != nil {
		return 0, err
	}

	err = s.storage.ChatMember().AddUserToChat(ctx, chatID, input.UserID)
	if err != nil {
		return 0, err
	}

	err = s.storage.ChatMember().AddUserToChat(ctx, chatID, input.FriendID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (s *chatService) CreateGroupChat(ctx context.Context, data *entity.CreateGroupChat) (*entity.Chat, error) {
	chatID, err := s.storage.Chat().Create(ctx, &entity.CreateChat{
		Type:       "group",
		Name:       data.Name,
		AvatarPath: data.AvatarPath,
	})
	if err != nil {
		return nil, err
	}

	err = s.storage.ChatMember().AddUserToChat(ctx, chatID, data.UserID)
	if err != nil {
		return nil, err
	}

	err = s.storage.ChatMember().AddUserToChat(ctx, chatID, data.FriendID)
	if err != nil {
		return nil, err
	}

	chat, err := s.storage.Chat().GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) UpdateGroupChat(ctx context.Context, userID, chatID uint, data *entity.UpdateChat) (*entity.Chat, error) {
	ok, err := s.storage.ChatMember().IsMemberOfChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Chat().Update(ctx, chatID, data)
	if err != nil {
		return nil, err
	}

	chat, err := s.storage.Chat().GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) DeletePrivateChat(ctx context.Context, userID, chatID uint) error {
	ok, err := s.storage.ChatMember().IsMemberOfChat(ctx, userID, chatID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrForbidden
	}

	chat, err := s.storage.Chat().GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Type != "private" {
		return domainErrors.ErrForbidden
	}

	return s.storage.Chat().Delete(ctx, chatID)
}
