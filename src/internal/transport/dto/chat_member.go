package dto

type AddUserToChatInput struct {
	UserID uint `json:"user_id" validate:"required"`
}
