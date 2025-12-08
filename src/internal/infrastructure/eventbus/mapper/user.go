package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewUser(user *entity.User) *dto.User {
	return &dto.User{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarPath,
	}
}
