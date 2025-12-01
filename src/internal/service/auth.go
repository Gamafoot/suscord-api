package service

import (
	"context"
	"suscord/internal/config"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/storage"
	"suscord/pkg/hash"

	pkgErrors "github.com/pkg/errors"
)

type authService struct {
	config  *config.Config
	storage storage.Storage
}

func NewAuthService(config *config.Config, storage storage.Storage) *authService {
	return &authService{
		config:  config,
		storage: storage,
	}
}

func (s *authService) Login(ctx context.Context, input *entity.LoginOrCreateInput) (string, error) {
	var (
		user *entity.User
	)

	hasher := hash.NewSHA1Hasher(s.config.Hash.Salt)

	user, err := s.storage.User().GetByUsername(ctx, input.Username)
	if err != nil {
		hash, err := hasher.Hash(input.Password)
		if err != nil {
			return "", pkgErrors.WithStack(err)
		}

		err = s.storage.User().Create(ctx, &entity.User{
			Username: input.Username,
			Password: hash,
		})
		if err != nil {
			return "", pkgErrors.WithStack(err)
		}

		user, err = s.storage.User().GetByUsername(ctx, input.Username)
		if err != nil {
			if pkgErrors.Is(err, domainErrors.ErrRecordNotFound) {
				return "", domainErrors.ErrInvalidLoginOrPassword
			}
			return "", pkgErrors.WithStack(err)
		}
	} else {
		hash, err := hasher.Hash(input.Password)
		if err != nil {
			return "", pkgErrors.WithStack(err)
		}

		if user.Password != hash {
			return "", domainErrors.ErrInvalidLoginOrPassword
		}
	}

	return s.createSession(ctx, user.ID)
}

func (s *authService) createSession(ctx context.Context, userID uint) (string, error) {
	var uuid string

	_, err := s.storage.Session().GetByUserID(ctx, userID)
	if err != nil {
		if pkgErrors.Is(err, domainErrors.ErrRecordNotFound) {
			uuid, err = s.storage.Session().Create(ctx, userID)
			if err != nil {
				return "", err
			}
		}
		return "", err
	} else {
		uuid, err = s.storage.Session().Update(ctx, userID)
		if err != nil {
			return "", err
		}
	}

	return uuid, nil
}
