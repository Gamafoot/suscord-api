package service

import (
	"suscord/internal/config"
	"suscord/internal/domain/broker"
	"suscord/internal/domain/logger"
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
	broker broker.Broker,
	logger logger.Logger,
) *_service {
	return &_service{
		user:       NewUserService(storage),
		auth:       NewAuthService(cfg, storage),
		chat:       NewChatService(cfg, storage, broker, logger),
		chatMember: NewChatMemberService(cfg, storage, broker, logger),
		message:    NewMessageService(cfg, storage, broker, logger),
		attachment: NewAttachmentService(cfg, storage, broker, logger),
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
