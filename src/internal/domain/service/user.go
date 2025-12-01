package service

import (
	"context"
	"suscord/internal/domain/entity"
)

type UserService interface {
	GetByID(ctx context.Context, userID uint) (*entity.User, error)
}
