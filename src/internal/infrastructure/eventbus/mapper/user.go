package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
	"suscord/pkg/urlpath"
)

func NewUser(user *entity.User, mediaURL string) *dto.User {
	return &dto.User{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: urlpath.GetMediaURL(mediaURL, user.AvatarPath),
	}
}
