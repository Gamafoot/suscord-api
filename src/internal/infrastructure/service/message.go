package service

import (
	"context"
	"mime"
	"mime/multipart"
	"path/filepath"
	"suscord/internal/config"
	"suscord/internal/domain/broker"
	brokerMsg "suscord/internal/domain/broker/message"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/logger"
	"suscord/internal/domain/storage"
)

type messageService struct {
	cfg     *config.Config
	storage storage.Storage
	broker  broker.Broker
	logger  logger.Logger
}

func NewMessageService(
	cfg *config.Config,
	storage storage.Storage,
	broker broker.Broker,
	logger logger.Logger,
) *messageService {
	return &messageService{
		cfg:     cfg,
		storage: storage,
		broker:  broker,
		logger:  logger,
	}
}

func (s *messageService) GetChatMessages(ctx context.Context, input *entity.GetMessagesInput) ([]*entity.Message, error) {
	ok, err := s.storage.Database().ChatMember().IsMemberOfChat(ctx, input.UserID, input.ChatID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	messages, err := s.storage.Database().Message().GetMessages(ctx, input.ChatID, input.LastMessageID, input.Limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *messageService) Create(ctx context.Context, userID, chatID uint, data *entity.CreateMessage, files []*multipart.FileHeader) (*entity.Message, error) {
	messageID, err := s.storage.Database().Message().Create(ctx, userID, chatID, data)
	if err != nil {
		return nil, err
	}

	message, err := s.storage.Database().Message().GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if len(files) > 0 {
		attachments, err := s.createAttachments(ctx, messageID, files)
		if err != nil {
			return nil, err
		}
		message.Attachments = attachments
	}

	err = s.broker.Publish(ctx, brokerMsg.NewMessageCreated(message, s.cfg.Media.Url))
	if err != nil {
		s.logger.Err(err, logger.Field{
			Key:   "message",
			Value: message,
		})
	}

	return message, nil
}

func (s *messageService) Update(ctx context.Context, userID, messageID uint, data *entity.UpdateMessage) (*entity.Message, error) {
	ok, err := s.storage.Database().Message().IsOwner(ctx, userID, messageID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Database().Message().Update(ctx, messageID, data)
	if err != nil {
		return nil, err
	}

	message, err := s.storage.Database().Message().GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	err = s.broker.Publish(ctx, brokerMsg.NewMessageUpdated(message, s.cfg.Media.Url))
	if err != nil {
		s.logger.Err(err, logger.Field{
			Key:   "message",
			Value: message,
		})
	}

	return message, nil
}

func (s *messageService) Delete(ctx context.Context, userID, messageID uint) error {
	ok, err := s.storage.Database().Message().IsOwner(ctx, userID, messageID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrUserIsNotMemberOfChat
	}

	message, err := s.storage.Database().Message().GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	err = s.storage.Database().Message().Delete(ctx, messageID)
	if err != nil {
		return err
	}

	err = s.broker.Publish(ctx, brokerMsg.NewMessageDeleted(message.ChatID, message.ID, userID))
	if err != nil {
		s.logger.Err(err,
			logger.Field{
				Key:   "chat_id",
				Value: message.ChatID,
			},
			logger.Field{
				Key:   "message_id",
				Value: message.ID,
			},
		)
	}

	return nil
}

func (s *messageService) createAttachments(
	ctx context.Context,
	messageID uint,
	files []*multipart.FileHeader,
) ([]*entity.Attachment, error) {
	attachments := make([]*entity.Attachment, len(files))

	for i, file := range files {
		mimetype := mime.TypeByExtension(filepath.Ext(file.Filename))

		filepath, err := s.storage.File().UploadFile(file, "messages")
		if err != nil {
			return nil, err
		}

		data := &entity.CreateAttachment{
			FilePath: filepath,
			FileSize: file.Size,
			MimeType: mimetype,
		}

		attachmentID, err := s.storage.Database().Attachment().Create(ctx, messageID, data)
		if err != nil {
			return nil, err
		}

		attachment, err := s.storage.Database().Attachment().GetByID(ctx, attachmentID)
		if err != nil {
			return nil, err
		}

		attachments[i] = attachment
	}

	return attachments, nil
}
