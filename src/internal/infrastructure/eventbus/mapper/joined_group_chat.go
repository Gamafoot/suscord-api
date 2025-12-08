package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewJoinedGroupChat(chat *entity.Chat, userID uint, dontSend ...bool) *dto.JoinedGroupChat {
	_dontSend := false

	if len(dontSend) > 0 {
		_dontSend = dontSend[0]
	}

	return &dto.JoinedGroupChat{
		Chat:     NewChat(chat),
		UserID:   userID,
		DontSend: _dontSend,
	}
}
