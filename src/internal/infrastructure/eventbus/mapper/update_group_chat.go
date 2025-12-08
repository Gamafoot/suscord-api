package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewUpdateGroupChat(chat *entity.Chat, exceptUserID uint) *dto.UpdateGroupChat {
	return &dto.UpdateGroupChat{
		Chat:         NewChat(chat),
		ExceptUserID: exceptUserID,
	}
}
