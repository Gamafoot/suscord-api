package service

import (
	"context"
	"suscord/internal/domain/entity"
)

type UserService interface {
	GetByID(ctx context.Context, userID uint) (*entity.User, error)
	SearchUsers(ctx context.Context, userID uint, username string) ([]*entity.User, error)
}
