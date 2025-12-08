package dto

type Chat struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	AvatarPath string `json:"avatar_path"`
	Type       string `json:"type"`
}
