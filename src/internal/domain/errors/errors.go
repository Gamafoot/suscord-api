package errors

import (
	"errors"
)

var (
	ErrIsNotDigit             = errors.New("is not digit")
	ErrIsNotPositiveDigit     = errors.New("digit must be positive")
	ErrLoginIsExists          = errors.New("login is exists")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrBadRequest             = errors.New("invalid request data")
	ErrRecordNotFound         = errors.New("record no found")
	ErrChatFull               = errors.New("the chat is full")
	ErrFailToJoinChat         = errors.New("fail to join chat")
	ErrUserIsNotMemberOfChat  = errors.New("user is not member of chat")
	ErrIsNotOwner             = errors.New("пользователь не является владельцем ресурса")
	ErrForbidden              = errors.New("пользователь не имеет доступа к ресурсу")
	ErrInvalidUserID          = errors.New("неверный user_id")

	ErrInvalidFileExtention = errors.New("запрещенный формат файла")
	ErrInvalidFile          = errors.New("нет доступа к файлу")
	ErrIsNotImage           = errors.New("файл не является изображением")
)
