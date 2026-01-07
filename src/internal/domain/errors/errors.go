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
	ErrIsNotOwner             = errors.New("user is not owner")
	ErrForbidden              = errors.New("no permissions")
	ErrInvalidUserID          = errors.New("invalid user_id")

	ErrInvalidFile = errors.New("invalid file")
	ErrIsNotImage  = errors.New("is not image")
)
