package storage

import (
	"suscord/internal/domain/storage"

	"gorm.io/gorm"
)

type _storage struct {
	user              *userStorage
	chat              *chatStorage
	chatMember        *chatMemberStorage
	message           *messageStorage
	messageAttachment *messageAttachmentStorage
	session           *sessionStorage
}

func NewGormStorage(db *gorm.DB) *_storage {
	return &_storage{
		user:              NewUserStorage(db),
		chat:              NewChatStorage(db),
		chatMember:        NewChatMemberStorage(db),
		message:           NewMessageStorage(db),
		messageAttachment: NewMessageAttachmentStorage(db),
		session:           NewSessionStorage(db),
	}
}

func (s *_storage) User() storage.UserStorage {
	return s.user
}

func (s *_storage) Chat() storage.ChatStorage {
	return s.chat
}

func (s *_storage) ChatMember() storage.ChatMemberStorage {
	return s.chatMember
}

func (s *_storage) Message() storage.MessageStorage {
	return s.message
}

func (s *_storage) MessageAttachment() storage.MessageAttachmentStorage {
	return s.messageAttachment
}

func (s *_storage) Session() storage.SessionStorage {
	return s.session
}
