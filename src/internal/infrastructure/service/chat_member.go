package service

import (
	"context"
	"strconv"
	"suscord/internal/config"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/eventbus/dto"
	"suscord/internal/domain/storage"
	"suscord/internal/infrastructure/eventbus/mapper"
	"time"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
)

type chatMemberService struct {
	cfg      *config.Config
	storage  storage.Storage
	eventbus eventbus.Bus
}

func NewChatMemberService(cfg *config.Config, storage storage.Storage, eventbus eventbus.Bus) *chatMemberService {
	return &chatMemberService{
		cfg:      cfg,
		storage:  storage,
		eventbus: eventbus,
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

	err = s.storage.Cache().Set(ctx, uuid.String(), chatID, 10*time.Second)
	if err != nil {
		return err
	}

	s.eventbus.Publish(&dto.InviteToRoom{
		Code:   uuid.String(),
		UserID: userID,
		ChatID: chatID,
	})

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

	data := mapper.NewUserInChat(uint(chatID), user, s.cfg.Media.Url)
	s.eventbus.Publish(data)

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

	data := mapper.NewLeftChat(chatID, userID)
	s.eventbus.Publish(data)

	return nil
}
