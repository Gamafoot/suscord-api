package model

type User struct {
	ID         uint
	Username   string `gorm:"type:varchar(20)"`
	Password   string `gorm:"type:text"`
	AvatarPath string `gorm:"type:varchar(255)"`
	FriendCode string `gorm:"type:varchar(20)"`
}

func (u User) TableName() string {
	return "users"
}
