package dto

type Chat struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_path"`
	Type      string `json:"type"`
}
