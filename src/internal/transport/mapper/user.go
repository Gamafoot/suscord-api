package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/transport/dto"
	"suscord/pkg/urlpath"
)

func NewUser(user *entity.User, mediaURL string) *dto.User {
	return &dto.User{
		ID:        user.ID,
		Username:  user.Username,
		AvatarUrl: urlpath.GetMediaURL(mediaURL, user.AvatarPath),
	}
}
