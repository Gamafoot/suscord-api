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

type chatStorage struct {
	db *gorm.DB
}

func NewChatStorage(db *gorm.DB) *chatStorage {
	return &chatStorage{db: db}
}

func (s *chatStorage) GetByID(ctx context.Context, chatID uint) (*entity.Chat, error) {
	chat := new(model.Chat)
	if err := s.db.WithContext(ctx).First(&chat, "id = ?", chatID).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return chatModelToDomain(chat), nil
}

func (s *chatStorage) GetUserChat(ctx context.Context, userID, chatID uint) (*entity.Chat, error) {
	chat := new(model.Chat)
	err := s.db.WithContext(ctx).Raw("select * from get_user_chat(?, ?)", chatID, userID).Scan(&chat).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return chatModelToDomain(chat), nil
}

func (s *chatStorage) GetUserChats(ctx context.Context, userID uint) ([]*entity.Chat, error) {
	chats := make([]*model.Chat, 0)
	err := s.db.WithContext(ctx).Raw("select * from get_user_chats(?)", userID).Scan(&chats).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	chatDomains := make([]*entity.Chat, len(chats))

	for i, chat := range chats {
		chatDomains[i] = chatModelToDomain(chat)
	}

	return chatDomains, nil
}

func (s *chatStorage) GetPrivateChatID(ctx context.Context, userID, friendID uint) (uint, error) {
	chat := new(entity.Chat)
	if err := s.db.WithContext(ctx).First(chat, "user_id = ? AND friend_id = ?", userID, friendID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domainErrors.ErrRecordNotFound
		}
		return 0, pkgErrors.WithStack(err)
	}
	return chat.ID, nil
}

func (s *chatStorage) Create(ctx context.Context, data *entity.CreateChat) (uint, error) {
	chat := &model.Chat{
		Type:       data.Type,
		Name:       data.Name,
		AvatarPath: data.AvatarPath,
	}

	if err := s.db.WithContext(ctx).Create(chat).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	return chat.ID, nil
}

func (s *chatStorage) Update(ctx context.Context, chatID uint, data *entity.UpdateChat) error {
	updateData := make(map[string]any)

	if data.Name != nil {
		updateData["name"] = data.Name
	}
	if data.AvatarPath != nil {
		updateData["avatar_path"] = data.AvatarPath
	}

	err := s.db.WithContext(ctx).Model(&model.Chat{}).Where("id = ?", chatID).Updates(data).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}

func (s *chatStorage) Delete(ctx context.Context, chatID uint) error {
	if err := s.db.WithContext(ctx).Delete(&entity.Chat{ID: chatID}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func chatModelToDomain(chat *model.Chat) *entity.Chat {
	return &entity.Chat{
		ID:         chat.ID,
		Name:       chat.Name,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
	}
}

func chatDomainToModel(chat *entity.Chat) *model.Chat {
	return &model.Chat{
		ID:         chat.ID,
		Name:       chat.Name,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
	}
}
