package storage

import (
	"context"
	"suscord/internal/database/gorm/model"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type messageStorage struct {
	db *gorm.DB
}

func NewMessageStorage(db *gorm.DB) *messageStorage {
	return &messageStorage{db: db}
}

func (s *messageStorage) GetMessages(ctx context.Context, chatID, lastMessageID uint, limit int) ([]*entity.Message, error) {
	messages := make([]*model.Message, 0)

	var (
		sql  string
		args []interface{}
	)

	if lastMessageID == 0 {
		sql = "select * from messages where chat_id = ? ORDER BY id DESC LIMIT ?"
		args = []interface{}{chatID, limit}
	} else {
		sql = "select * from messages where chat_id = ? AND id < ? ORDER BY id DESC LIMIT ?"
		args = []interface{}{chatID, lastMessageID, limit}
	}

	if err := s.db.WithContext(ctx).Raw(sql, args...).Scan(&messages).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	messageDomains := make([]*entity.Message, len(messages))
	for i, message := range messages {
		messageDomains[i] = messageModelToDomain(message)
	}

	return messageDomains, nil
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

func (s *messageStorage) Create(ctx context.Context, userID, chatID uint, data *entity.CreateMessage) (uint, error) {
	message := &model.Message{
		UserID:  userID,
		ChatID:  chatID,
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
	return &entity.Message{
		ID:        message.ID,
		ChatID:    message.ChatID,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}
