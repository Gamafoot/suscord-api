package service

import (
	"suscord/internal/config"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/service"
	"suscord/internal/domain/storage"
)

type _service struct {
	user       *userService
	auth       *authService
	chat       *chatService
	chatMember *chatMemberService
	message    *messageService
	attachment *attachmentService
	file       *fileService
}

func NewService(
	cfg *config.Config,
	storage storage.Storage,
	eventbus eventbus.Bus,
) *_service {
	return &_service{
		user:       NewUserService(storage),
		auth:       NewAuthService(cfg, storage),
		chat:       NewChatService(storage, eventbus),
		chatMember: NewChatMemberService(storage, eventbus),
		message:    NewMessageService(cfg, storage, eventbus),
		attachment: NewAttachmentService(cfg, storage, eventbus),
		file:       NewFileService(storage),
	}
}

func (s *_service) User() service.UserService {
	return s.user
}

func (s *_service) Auth() service.AuthService {
	return s.auth
}

func (s *_service) Chat() service.ChatService {
	return s.chat
}

func (s *_service) ChatMember() service.ChatMemberService {
	return s.chatMember
}

func (s *_service) Message() service.MessageService {
	return s.message
}

func (s *_service) Attachment() service.AttachementService {
	return s.attachment
}

func (s *_service) File() service.FileService {
	return s.file
}
