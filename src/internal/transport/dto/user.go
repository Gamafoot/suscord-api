package dto

type User struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type Me struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	AvatarUrl  string `json:"avatar_url"`
	FriendCode string `json:"friend_code"`
}
