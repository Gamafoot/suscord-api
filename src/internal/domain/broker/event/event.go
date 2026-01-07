package event

const (
	OnMessageCreated        = "chat.message.created"
	OnMessageUpdated        = "chat.message.updated"
	OnMessageDeleted        = "chat.message.deleted"
	OnChatUpdated           = "chat.updated"
	OnChatDeleted           = "chat.deleted"
	OnUserInvited           = "chat.user.invited"
	OnUserJoinedGroupChat   = "chat.group.user.joined"
	OnUserJoinedPrivateChat = "chat.private.user.joined"
	OnUserLeft              = "chat.user.leave"
)
