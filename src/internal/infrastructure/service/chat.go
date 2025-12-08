package service

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
	"suscord/internal/infrastructure/eventbus/mapper"

	"github.com/pkg/errors"
)

type chatService struct {
	storage  storage.Storage
	eventbus eventbus.Bus
}

func NewChatService(storage storage.Storage, eventbus eventbus.Bus) *chatService {
	return &chatService{
		storage:  storage,
		eventbus: eventbus,
	}
}

func (s *chatService) GetUserChats(ctx context.Context, userID uint, searchPattern string) ([]*entity.Chat, error) {
	if searchPattern != "" {
		return s.storage.Database().Chat().SearchUserChats(ctx, userID, searchPattern)
	}
	return s.storage.Database().Chat().GetUserChats(ctx, userID)
}

func (s *chatService) GetOrCreatePrivateChat(ctx context.Context, input *entity.CreatePrivateChat) (*entity.Chat, error) {
	createChat := false

	chatID, err := s.storage.Database().ChatMember().GetPrivateChatID(ctx, input.UserID, input.FriendID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			chatID, err = s.createPrivateChat(ctx, input)
			if err != nil {
				return nil, err
			}

			createChat = true
		} else {
			return nil, err
		}
	}

	chat, err := s.storage.Database().Chat().GetUserChat(ctx, chatID, input.UserID)
	if err != nil {
		return nil, err
	}

	if createChat {
		data := mapper.NewJoinedPrivateChat(chat, input.UserID, true)
		s.eventbus.Publish(data)

		chatFriend, err := s.storage.Database().Chat().GetUserChat(ctx, chatID, input.FriendID)
		if err != nil {
			return nil, err
		}

		data = mapper.NewJoinedPrivateChat(chatFriend, input.FriendID)
		s.eventbus.Publish(data)
	}

	return chat, nil
}

func (s *chatService) createPrivateChat(ctx context.Context, input *entity.CreatePrivateChat) (uint, error) {
	chatID, err := s.storage.Database().Chat().Create(ctx, &entity.CreateChat{Type: "private"})
	if err != nil {
		return 0, err
	}

	err = s.storage.Database().ChatMember().AddUserToChat(ctx, input.UserID, chatID)
	if err != nil {
		return 0, err
	}

	err = s.storage.Database().ChatMember().AddUserToChat(ctx, input.FriendID, chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (s *chatService) CreateGroupChat(ctx context.Context, userID uint, data *entity.CreateGroupChat) (*entity.Chat, error) {
	chatID, err := s.storage.Database().Chat().Create(ctx, &entity.CreateChat{
		Type:       "group",
		Name:       data.Name,
		AvatarPath: data.AvatarPath,
	})
	if err != nil {
		return nil, err
	}

	err = s.storage.Database().ChatMember().AddUserToChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	chat, err := s.storage.Database().Chat().GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	event := mapper.NewJoinedGroupChat(chat, userID, true)
	s.eventbus.Publish(event)

	return chat, nil
}

func (s *chatService) UpdateGroupChat(
	ctx context.Context,
	userID uint,
	chatID uint,
	data *entity.UpdateChat,
) (*entity.Chat, error) {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Database().Chat().Update(ctx, chatID, data)
	if err != nil {
		return nil, err
	}

	chat, err := s.storage.Database().Chat().GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	s.eventbus.Publish(mapper.NewUpdateGroupChat(chat, userID))

	return chat, nil
}

func (s *chatService) DeletePrivateChat(ctx context.Context, userID, chatID uint) error {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, userID, chatID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrForbidden
	}

	chat, err := s.storage.Database().Chat().GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Type != "private" {
		return domainErrors.ErrForbidden
	}

	data := mapper.NewDeleteChat(chatID, userID)
	s.eventbus.Publish(data)

	return s.storage.Database().Chat().Delete(ctx, chatID)
}
