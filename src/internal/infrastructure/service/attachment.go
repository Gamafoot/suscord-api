package service

import (
	"context"
	"suscord/internal/config"
	"suscord/internal/domain/broker"
	"suscord/internal/domain/broker/event"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/logger"
	"suscord/internal/domain/storage"
)

type attachmentService struct {
	cfg     *config.Config
	storage storage.Storage
	broker  broker.Broker
	logger  logger.Logger
}

func NewAttachmentService(
	cfg *config.Config,
	storage storage.Storage,
	broker broker.Broker,
	logger logger.Logger,
) *attachmentService {
	return &attachmentService{
		cfg:     cfg,
		storage: storage,
		broker:  broker,
		logger:  logger,
	}
}

func (s *attachmentService) Delete(ctx context.Context, userID, attachmentID uint) error {
	ok, err := s.storage.Database().Attachment().IsOwner(ctx, userID, attachmentID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrUserIsNotMemberOfChat
	}

	message, err := s.storage.Database().Message().GetByAttachmentID(ctx, attachmentID)
	if err != nil {
		return err
	}

	err = s.storage.Database().Attachment().Delete(ctx, attachmentID)
	if err != nil {
		return err
	}

	brokerCtx, cansel := context.WithTimeout(context.Background(), s.cfg.Broker.Timeout)
	defer cansel()

	fields := []logger.Field{
		logger.Field{
			Key:   "chat_id",
			Value: message.ChatID,
		},
		logger.Field{
			Key:   "message_id",
			Value: message.ID,
		},
	}

	err = s.broker.Publish(brokerCtx, event.NewMessageDeleted(message.ChatID, message.ID))
	if err != nil {
		s.logger.Err(err, fields...)
	}

	return nil
}
