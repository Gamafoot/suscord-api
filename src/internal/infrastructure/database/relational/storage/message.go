package storage

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/infrastructure/database/relational/model"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type messageStorage struct {
	db *gorm.DB
}

func NewMessageStorage(db *gorm.DB) *messageStorage {
	return &messageStorage{db: db}
}

func (s *messageStorage) GetMessages(
	ctx context.Context,
	chatID, lastMessageID uint,
	limit int,
) ([]*entity.Message, error) {

	messages := make([]*model.Message, 0)

	db := s.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("chat_id = ?", chatID)

	if lastMessageID != 0 {
		db = db.Where("id < ?", lastMessageID)
	}

	if err := db.
		Order("id DESC").
		Limit(limit).
		Preload("Attachments").
		Find(&messages).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	// конвертация в domain
	result := make([]*entity.Message, len(messages))
	for i, m := range messages {
		result[i] = messageModelToDomain(m)
	}

	return result, nil
}

func (s *messageStorage) GetByID(ctx context.Context, messageID uint) (*entity.Message, error) {
	message := new(model.Message)
	if err := s.db.WithContext(ctx).First(&message, "id = ?", messageID).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}
	return messageModelToDomain(message), nil
}

func (s *messageStorage) GetByAttachmentID(ctx context.Context, attachmentID uint) (*entity.Message, error) {
	message := new(model.Message)

	err := s.db.WithContext(ctx).
		Joins("INNER JOIN message_attachments a ON a.message_id = messages.id").
		First(&message, "a.id = ?", attachmentID).
		Error

	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}

	return messageModelToDomain(message), nil
}

func (s *messageStorage) Create(ctx context.Context, userID, chatID uint, data *entity.CreateMessage) (uint, error) {
	message := &model.Message{
		UserID:  userID,
		ChatID:  chatID,
		Type:    data.Type,
		Content: data.Content,
	}
	if err := s.db.WithContext(ctx).Create(message).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}
	return message.ID, nil
}

func (s *messageStorage) Update(ctx context.Context, messageID uint, data *entity.UpdateMessage) error {
	message := &model.Message{
		ID:      messageID,
		Content: data.Content,
	}
	err := s.db.WithContext(ctx).Updates(message).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *messageStorage) Delete(ctx context.Context, messageID uint) error {
	if err := s.db.WithContext(ctx).Delete(&entity.Message{ID: messageID}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *messageStorage) IsOwner(ctx context.Context, userID, messageID uint) (bool, error) {
	message := new(model.Message)
	if err := s.db.WithContext(ctx).First(&message, "id = ? AND user_id = ?", messageID, userID).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return false, domainErrors.ErrRecordNotFound
		}
		return false, pkgErrors.WithStack(err)
	}
	return true, nil
}

func messageModelToDomain(message *model.Message) *entity.Message {
	attachments := make([]*entity.Attachment, len(message.Attachments))
	for i, attachment := range message.Attachments {
		attachments[i] = &entity.Attachment{
			ID:        attachment.ID,
			MessageID: attachment.MessageID,
			FilePath:  attachment.FilePath,
			FileSize:  attachment.FileSize,
			MimeType:  attachment.MimeType,
		}
	}

	return &entity.Message{
		ID:          message.ID,
		ChatID:      message.ChatID,
		UserID:      message.UserID,
		Type:        message.Type,
		Content:     message.Content,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
		Attachments: attachments,
	}
}
