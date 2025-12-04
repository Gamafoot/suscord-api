package model

type Chat struct {
	ID         uint
	Name       string `gorm:"type:varchar(50)"`
	AvatarPath string `gorm:"type:varchar(255)"`
	Type       string `gorm:"type:varchar(20)"`
}

func (c Chat) TableName() string {
	return "chats"
}
