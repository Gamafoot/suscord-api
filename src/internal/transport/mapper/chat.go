package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/transport/dto"
	"suscord/pkg/urlpath"
)

func NewChat(chat *entity.Chat, mediaURL string) *dto.Chat {
	return &dto.Chat{
		ID:        chat.ID,
		Type:      chat.Type,
		Name:      chat.Name,
		AvatarUrl: urlpath.GetMediaURL(mediaURL, chat.AvatarPath),
	}
}
