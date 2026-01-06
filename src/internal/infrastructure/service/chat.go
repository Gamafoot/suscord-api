package service

import (
	"context"
	"suscord/internal/config"
	"suscord/internal/domain/broker"
	brokerMsg "suscord/internal/domain/broker/message"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/logger"
	"suscord/internal/domain/storage"

	"github.com/pkg/errors"
)

type chatService struct {
	cfg     *config.Config
	storage storage.Storage
	broker  broker.Broker
	logger  logger.Logger
}

func NewChatService(
	cfg *config.Config,
	storage storage.Storage,
	broker broker.Broker,
	logger logger.Logger,
) *chatService {
	return &chatService{
		cfg:     cfg,
		storage: storage,
		broker:  broker,
		logger:  logger,
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
		user, err := s.storage.Database().User().GetByID(ctx, input.UserID)
		if err != nil {
			return nil, err
		}

		err = s.broker.Publish(ctx, brokerMsg.NewUserJoinedPrivateChat(chatID, user, s.cfg.Static.URL))
		if err != nil {
			s.logger.Err(err,
				logger.Field{
					Key:   "chat_id",
					Value: chatID,
				},
				logger.Field{
					Key:   "user_id",
					Value: input.UserID,
				},
			)
		}

		user, err = s.storage.Database().User().GetByID(ctx, input.FriendID)
		if err != nil {
			return nil, err
		}

		err = s.broker.Publish(ctx, brokerMsg.NewUserJoinedPrivateChat(chatID, user, s.cfg.Static.URL))
		if err != nil {
			s.logger.Err(err,
				logger.Field{
					Key:   "chat_id",
					Value: chatID,
				},
				logger.Field{
					Key:   "user_id",
					Value: input.FriendID,
				},
			)
		}
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

	err = s.broker.Publish(ctx, brokerMsg.NewChatUpdated(chat, s.cfg.Media.Url))
	if err != nil {
		s.logger.Err(err,
			logger.Field{
				Key:   "chat_id",
				Value: chatID,
			},
		)
	}

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

	err = s.broker.Publish(ctx, brokerMsg.NewChatDeleted(chatID, chat.Name))
	if err != nil {
		s.logger.Err(err,
			logger.Field{
				Key:   "chat_id",
				Value: chatID,
			},
		)
	}

	return s.storage.Database().Chat().Delete(ctx, chatID)
}
