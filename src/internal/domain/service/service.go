package service

type Service interface {
	User() UserService
	Auth() AuthService
	Chat() ChatService
	ChatMember() ChatMemberService
	Message() MessageService
	Attachment() AttachementService
	File() FileService
}
