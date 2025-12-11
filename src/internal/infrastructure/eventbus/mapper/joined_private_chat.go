package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewJoinedPrivateChat(chat *entity.Chat, userID uint, mediaURL string, dontSend ...bool) *dto.JoinedPrivateChat {
	_dontSend := false

	if len(dontSend) > 0 {
		_dontSend = dontSend[0]
	}

	return &dto.JoinedPrivateChat{
		Chat:     NewChat(chat, mediaURL),
		UserID:   userID,
		DontSend: _dontSend,
	}
}
