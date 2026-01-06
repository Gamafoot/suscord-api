package service

import (
	"context"
	"strconv"
	"suscord/internal/config"
	"suscord/internal/domain/broker"
	brokerMsg "suscord/internal/domain/broker/message"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/logger"
	"suscord/internal/domain/storage"
	"time"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
)

type chatMemberService struct {
	cfg     *config.Config
	storage storage.Storage
	broker  broker.Broker
	logger  logger.Logger
}

func NewChatMemberService(
	cfg *config.Config,
	storage storage.Storage,
	eventbus broker.Broker,
	logger logger.Logger,
) *chatMemberService {
	return &chatMemberService{
		cfg:     cfg,
		storage: storage,
		broker:  eventbus,
		logger:  logger,
	}
}

func (s *chatMemberService) GetNonMembers(ctx context.Context, chatID, memberID uint) ([]*entity.User, error) {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, memberID, chatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	users, err := s.storage.Database().ChatMember().GetChatMembers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *chatMemberService) GetNotChatMembers(ctx context.Context, chatID, memberID uint) ([]*entity.User, error) {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, memberID, chatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	users, err := s.storage.Database().ChatMember().GetNonMembers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *chatMemberService) SendInvite(ctx context.Context, ownerID, chatID, userID uint) error {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, ownerID, chatID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrUserIsNotMemberOfChat
	}

	chat, err := s.storage.Database().Chat().GetByID(ctx, chatID)
	if err != nil {
		return err
	}

	if chat.Type != "group" {
		return domainErrors.ErrChatIsNotGroup
	}

	uuid := uuid.New()
	code := uuid.String()

	err = s.storage.Cache().Set(ctx, code, chatID, 10*time.Second)
	if err != nil {
		return err
	}

	err = s.broker.Publish(ctx, brokerMsg.NewUserInvited(chatID, userID, code))
	if err != nil {
		s.logger.Err(err,
			logger.Field{
				Key:   "chat_id",
				Value: chatID,
			},
			logger.Field{
				Key:   "user_id",
				Value: userID,
			},
		)
	}

	return nil
}

func (s *chatMemberService) AcceptInvite(ctx context.Context, userID uint, code string) error {
	value, err := s.storage.Cache().Get(ctx, code)
	if err != nil {
		return err
	}

	chatID, err := strconv.Atoi(value)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	err = s.storage.Database().ChatMember().AddUserToChat(ctx, userID, uint(chatID))
	if err != nil {
		return err
	}

	user, err := s.storage.Database().User().GetByID(ctx, userID)
	if err != nil {
		return err
	}

	err = s.broker.Publish(ctx, brokerMsg.NewUserJoinedGroupChat(uint(chatID), user, s.cfg.Media.Url))
	if err != nil {
		s.logger.Err(err,
			logger.Field{
				Key:   "chat_id",
				Value: chatID,
			},
			logger.Field{
				Key:   "user_id",
				Value: userID,
			},
		)
	}

	return nil
}

func (s *chatMemberService) LeaveFromChat(ctx context.Context, chatID, userID uint) error {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, userID, chatID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Database().ChatMember().Delete(ctx, userID, chatID)
	if err != nil {
		return err
	}

	err = s.broker.Publish(ctx, brokerMsg.NewUserLeft(chatID, userID))
	if err != nil {
		s.logger.Err(err,
			logger.Field{
				Key:   "chat_id",
				Value: chatID,
			},
			logger.Field{
				Key:   "user_id",
				Value: userID,
			},
		)
	}

	return nil
}
