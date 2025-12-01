package service

import (
	"context"
	"suscord/internal/domain/entity"
	"suscord/internal/domain/storage"

	pkgErrors "github.com/pkg/errors"
)

type userService struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) *userService {
	return &userService{
		storage: storage,
	}
}

func (s *userService) GetByID(ctx context.Context, userID uint) (*entity.User, error) {
	user, err := s.storage.User().GetByID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get user")
	}
	return user, nil
}
