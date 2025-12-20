package dto

type Client struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}
