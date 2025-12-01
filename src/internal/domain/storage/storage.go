package storage

type Storage interface {
	User() UserStorage
	Chat() ChatStorage
	ChatMember() ChatMemberStorage
	Message() MessageStorage
	MessageAttachment() MessageAttachmentStorage
	Session() SessionStorage
}
