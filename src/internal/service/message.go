package service

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
)

type messageService struct {
	storage  storage.Storage
	eventbus eventbus.Bus
}

func NewMessageService(storage storage.Storage, eventbus eventbus.Bus) *messageService {
	return &messageService{
		storage:  storage,
		eventbus: eventbus,
	}
}

func (s *messageService) GetChatMessages(ctx context.Context, input *entity.GetMessagesInput) ([]*entity.Message, error) {
	ok, err := s.storage.ChatMember().IsMemberOfChat(ctx, input.UserID, input.ChatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	messages, err := s.storage.Message().GetMessages(ctx, input.ChatID, input.LastMessageID, input.Limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *messageService) Create(ctx context.Context, userID, chatID uint, input *entity.CreateMessage) (*entity.Message, error) {
	messageID, err := s.storage.Message().Create(ctx,
		userID,
		chatID,
		&entity.CreateMessage{
			Content: input.Content,
		},
	)
	if err != nil {
		return nil, err
	}

	message, err := s.storage.Message().GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	go s.eventbus.Publish(eventbus.MessageCreated{
		ID:        message.ID,
		ChatID:    message.ChatID,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	})

	return message, nil
}

func (s *messageService) Update(ctx context.Context, userID, messageID uint, data *entity.UpdateMessage) (*entity.Message, error) {
	ok, err := s.storage.Message().IsOwner(ctx, userID, messageID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Message().Update(ctx, messageID, data)
	if err != nil {
		return nil, err
	}

	return s.storage.Message().GetByID(ctx, messageID)
}

func (s *messageService) Delete(ctx context.Context, userID, messageID uint) error {
	ok, err := s.storage.Message().IsOwner(ctx, userID, messageID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Message().Delete(ctx, messageID)
	if err != nil {
		return err
	}

	return nil
}
