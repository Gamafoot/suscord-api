package service

import (
	"context"
	"suscord/internal/domain/entity"
)

type AuthService interface {
	Login(ctx context.Context, input *entity.LoginOrCreateInput) (string, error)
}
