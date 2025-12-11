package dto

type Chat struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type CreatePrivateChatRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

type CreateGroupChatRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateGroupChatInput struct {
	Name *string `form:"name"`
}
