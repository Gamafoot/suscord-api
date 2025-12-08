package entity

type Chat struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	AvatarPath string `json:"avatar_path"`
	Type       string `json:"type"`
}

type CreatePrivateChat struct {
	UserID   uint
	FriendID uint
}

type CreateGroupChat struct {
	Name       string
	AvatarPath string
}

type CreateChat struct {
	Type       string
	Name       string
	AvatarPath string
}

type UpdateChat struct {
	Name       *string
	AvatarPath *string
}
