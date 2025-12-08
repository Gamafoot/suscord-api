package dto

type User struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	AvatarPath string `json:"avatar_path"`
}

type Me struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	AvatarPath string `json:"avatar_path"`
	FriendCode string `json:"friend_code"`
}
