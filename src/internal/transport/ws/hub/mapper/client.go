package mapper

import (
	"suscord/internal/transport/ws/hub/dto"
	"suscord/internal/transport/ws/hub/model"
	"suscord/pkg/urlpath"
)

func NewClient(client *model.Client, mediaURL string) *dto.Client {
	if client == nil {
		return nil
	}

	return &dto.Client{
		ID:        client.ID,
		Username:  client.Username,
		AvatarUrl: urlpath.GetMediaURL(mediaURL, client.AvatarPath),
	}
}
