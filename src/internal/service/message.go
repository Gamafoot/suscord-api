package service

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"suscord/internal/config"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/eventbus/mapper"
	"suscord/internal/domain/storage"
	"time"

	pkgErrors "github.com/pkg/errors"
)

type messageService struct {
	cfg      *config.Config
	storage  storage.Storage
	eventbus eventbus.Bus
}

func NewMessageService(cfg *config.Config, storage storage.Storage, eventbus eventbus.Bus) *messageService {
	return &messageService{
		cfg:      cfg,
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

func (s *messageService) Create(ctx context.Context, userID, chatID uint, data *entity.CreateMessage, files []*multipart.FileHeader) (*entity.Message, error) {
	messageID, err := s.storage.Message().Create(ctx, userID, chatID, data)
	if err != nil {
		return nil, err
	}

	message, err := s.storage.Message().GetByID(ctx, messageID)
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

	go s.eventbus.Publish(mapper.NewCreatedMessage(message))

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

func (s *messageService) createAttachments(
	ctx context.Context,
	messageID uint,
	files []*multipart.FileHeader,
) ([]*entity.Attachment, error) {
	attachments := make([]*entity.Attachment, len(files))

	for i, file := range files {
		mimetype := mime.TypeByExtension(filepath.Ext(file.Filename))

		filepath, err := s.saveFile(file)
		if err != nil {
			return nil, err
		}

		data := &entity.CreateAttachment{
			FilePath: filepath,
			FileSize: file.Size,
			MimeType: mimetype,
		}

		attachmentID, err := s.storage.Attachment().Create(ctx, messageID, data)
		if err != nil {
			return nil, err
		}

		attachment, err := s.storage.Attachment().GetByID(ctx, attachmentID)
		if err != nil {
			return nil, err
		}

		attachments[i] = attachment
	}

	return attachments, nil
}

func (s *messageService) saveFile(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d"+ext, file.Filename, time.Now().UnixNano())

	var (
		rootpath string
		year     int
		month    int
	)

	if year == 0 && month == 0 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	rootpath = fmt.Sprintf("%s/%d/%d", s.cfg.Media.RootFolder, year, month)
	filepath := fmt.Sprintf("%s/%s", rootpath, filename)

	os.MkdirAll(rootpath, os.ModePerm)

	src, err := file.Open()
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}
	defer dst.Close()

	io.Copy(dst, src)

	return filepath, nil
}
