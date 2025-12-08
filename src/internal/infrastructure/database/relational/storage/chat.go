package storage

import (
	"context"
	"suscord/internal/domain/entity"
	"suscord/internal/infrastructure/database/relational/model"

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

func (s *chatStorage) GetUserChat(ctx context.Context, chatID, userID uint) (*entity.Chat, error) {
	chat := &model.Chat{}

	err := s.db.WithContext(ctx).Raw("select * from get_user_chat(?, ?)", chatID, userID).Scan(chat).Error
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

func (s *chatStorage) SearchUserChats(ctx context.Context, userID uint, searchPattern string) ([]*entity.Chat, error) {
	chats := make([]*model.Chat, 0)

	err := s.db.WithContext(ctx).
		Raw("select * from search_user_chats(?, ?)", userID, searchPattern).
		Scan(&chats).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	chatDomains := make([]*entity.Chat, len(chats))

	for i, chat := range chats {
		chatDomains[i] = chatModelToDomain(chat)
	}

	return chatDomains, nil
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
