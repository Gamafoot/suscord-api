package dto

type Chat struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	AvatarPath string `json:"avatar_path"`
}

type CreatePrivateChatRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

type CreateGroupChatRequest struct {
	Name       string `json:"name" validate:"required"`
	AvatarPath string `json:"avatar_path"`
}

type UpdateGroupChatInput struct {
	Name *string `form:"name"`
}
