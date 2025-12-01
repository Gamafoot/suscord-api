package model

type Friend struct {
	ID       uint
	UserID   uint
	FriendID uint

	User   User
	Friend User
}

func (f *Friend) TableName() string {
	return "friends"
}
