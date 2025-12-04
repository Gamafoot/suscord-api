package storage

import (
	"context"
	"suscord/internal/database/relational/model"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"

	"github.com/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type chatMemberStorage struct {
	db *gorm.DB
}

func NewChatMemberStorage(db *gorm.DB) *chatMemberStorage {
	return &chatMemberStorage{db: db}
}

func (s *chatMemberStorage) GetChatMembers(ctx context.Context, chatID uint) ([]*entity.User, error) {
	users := make([]*model.User, 0)

	sql := `
	SELECT *
	FROM users
	WHERE id IN (
		SELECT user_id 
		FROM chat_members 
		WHERE chat_id = ?
	);`

	if err := s.db.WithContext(ctx).Raw(sql, chatID).Scan(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}

	domainUsers := make([]*entity.User, len(users))
	for i, user := range users {
		domainUsers[i] = userModelToDomain(user)
	}

	return domainUsers, nil
}

func (s *chatMemberStorage) AddUserToChat(ctx context.Context, userID, chatID uint) error {
	member := &entity.ChatMember{
		ChatID: chatID,
		UserID: userID,
	}
	if err := s.db.WithContext(ctx).Create(member).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *chatMemberStorage) IsMemberOfChat(ctx context.Context, userID, chatID uint) (bool, error) {
	if err := s.db.WithContext(ctx).First(&entity.ChatMember{}, "chat_id = ? AND user_id = ?", chatID, userID).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, pkgErrors.WithStack(err)
	}
	return true, nil
}

func (s *chatMemberStorage) Delete(ctx context.Context, userID, chatID uint) error {
	if err := s.db.WithContext(ctx).Where("chat_id = ? AND user_id = ?", chatID, userID).Delete(&entity.ChatMember{}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}
