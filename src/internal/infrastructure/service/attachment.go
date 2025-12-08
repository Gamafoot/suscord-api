package service

import (
	"context"
	"suscord/internal/config"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/eventbus/events"
	"suscord/internal/domain/storage"
	"suscord/internal/infrastructure/eventbus/mapper"
)

type attachmentService struct {
	cfg      *config.Config
	storage  storage.Storage
	eventbus eventbus.Bus
}

func NewAttachmentService(cfg *config.Config, storage storage.Storage, eventbus eventbus.Bus) *attachmentService {
	return &attachmentService{
		cfg:      cfg,
		storage:  storage,
		eventbus: eventbus,
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

	data := mapper.NewMessage(events.EventMessageUpdate, message)
	s.eventbus.Publish(data)

	return nil
}
