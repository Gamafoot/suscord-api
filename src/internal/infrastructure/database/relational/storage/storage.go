package storage

import (
	"suscord/internal/domain/storage/database"

	"gorm.io/gorm"
)

type _storage struct {
	user       *userStorage
	chat       *chatStorage
	chatMember *chatMemberStorage
	message    *messageStorage
	attachment *attachmentStorage
	session    *sessionStorage
}

func NewGormStorage(db *gorm.DB) *_storage {
	return &_storage{
		user:       NewUserStorage(db),
		chat:       NewChatStorage(db),
		chatMember: NewChatMemberStorage(db),
		message:    NewMessageStorage(db),
		attachment: NewAttachmentStorage(db),
		session:    NewSessionStorage(db),
	}
}

func (s *_storage) User() database.UserStorage {
	return s.user
}

func (s *_storage) Chat() database.ChatStorage {
	return s.chat
}

func (s *_storage) ChatMember() database.ChatMemberStorage {
	return s.chatMember
}

func (s *_storage) Message() database.MessageStorage {
	return s.message
}

func (s *_storage) Attachment() database.AttachmentStorage {
	return s.attachment
}

func (s *_storage) Session() database.SessionStorage {
	return s.session
}
