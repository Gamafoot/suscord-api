package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/transport/dto"
)

func NewUser(user *entity.User) *dto.User {
	return &dto.User{
		ID:         user.ID,
		Username:   user.Username,
		AvatarPath: user.AvatarPath,
	}
}
