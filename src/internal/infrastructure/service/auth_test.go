package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"suscord/internal/config"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	dbmocks "suscord/internal/domain/storage/database/mocks"
	storagemocks "suscord/internal/domain/storage/mocks"
	"suscord/pkg/hash"

	pkgErrors "github.com/pkg/errors"
)

func TestAuthService_ExistingUserValidPassword(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Hash.Salt = "salt"

	username := "u"
	password := "pass"
	userID := uint(1)
	sessionUUID := "uuid-1"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)
	sessionRepo := dbmocks.NewMockSessionStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()
	db.EXPECT().Session().Return(sessionRepo).Maybe()

	hasher := hash.NewSHA1Hasher(cfg.Hash.Salt)
	passHash, err := hasher.Hash(password)
	require.NoError(t, err)

	userRepo.EXPECT().
		GetByUsername(mock.Anything, username).
		Return(&entity.User{ID: userID, Username: username, Password: passHash}, nil).
		Once()

	sessionRepo.EXPECT().
		GetByUserID(mock.Anything, userID).
		Return(nil, domainErrors.ErrRecordNotFound).
		Once()

	sessionRepo.EXPECT().
		Create(mock.Anything, userID).
		Return(sessionUUID, nil).
		Once()

	s := NewAuthService(cfg, storage)

	uuid, err := s.Login(ctx, &entity.LoginOrCreateInput{Username: username, Password: password})
	require.NoError(t, err)
	assert.Equal(t, sessionUUID, uuid)
}

func TestAuthService_Login_ExistingUserInvalidPassword(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Hash.Salt = "salt"

	username := "u"
	password := "pass"
	userID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	userRepo.EXPECT().
		GetByUsername(mock.Anything, username).
		Return(&entity.User{ID: userID, Username: username, Password: "some-other-hash"}, nil).
		Once()

	s := NewAuthService(cfg, storage)

	uuid, err := s.Login(ctx, &entity.LoginOrCreateInput{Username: username, Password: password})
	assert.ErrorIs(t, err, domainErrors.ErrInvalidLoginOrPassword)
	assert.Empty(t, uuid)
}

func TestAuthService_Login_NewUser_CreatesUserAndUpdatesSession(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Hash.Salt = "salt"

	username := "new"
	password := "pass"
	userID := uint(1)
	sessionUUID := "uuid-1"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)
	sessionRepo := dbmocks.NewMockSessionStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()
	db.EXPECT().Session().Return(sessionRepo).Maybe()

	hasher := hash.NewSHA1Hasher(cfg.Hash.Salt)
	passHash, err := hasher.Hash(password)
	require.NoError(t, err)

	userRepo.EXPECT().
		GetByUsername(mock.Anything, username).
		Return(nil, pkgErrors.WithStack(domainErrors.ErrRecordNotFound)).
		Once()

	userRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*entity.User")).
		Run(func(_ context.Context, u *entity.User) {
			assert.Equal(t, username, u.Username)
			assert.Equal(t, passHash, u.Password)
		}).
		Return(nil).
		Once()

	userRepo.EXPECT().
		GetByUsername(mock.Anything, username).
		Return(&entity.User{ID: userID, Username: username, Password: passHash}, nil).
		Once()

	sessionRepo.EXPECT().
		GetByUserID(mock.Anything, userID).
		Return(&entity.Session{UserID: userID}, nil).
		Once()

	sessionRepo.EXPECT().
		Update(mock.Anything, userID).
		Return(sessionUUID, nil).
		Once()

	s := NewAuthService(cfg, storage)

	uuid, err := s.Login(ctx, &entity.LoginOrCreateInput{Username: username, Password: password})
	require.NoError(t, err)
	assert.Equal(t, sessionUUID, uuid)
}

func TestAuthService_ReturnsInvalidLoginOrPassword(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Hash.Salt = "salt"

	username := "new"
	password := "pass"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	userRepo.EXPECT().
		GetByUsername(mock.Anything, username).
		Return(nil, pkgErrors.WithStack(domainErrors.ErrRecordNotFound)).
		Once()

	userRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*entity.User")).
		Return(nil).
		Once()

	userRepo.EXPECT().
		GetByUsername(mock.Anything, username).
		Return(nil, domainErrors.ErrRecordNotFound).
		Once()

	s := NewAuthService(cfg, storage)

	uuid, err := s.Login(ctx, &entity.LoginOrCreateInput{Username: username, Password: password})
	assert.ErrorIs(t, err, domainErrors.ErrInvalidLoginOrPassword)
	assert.Empty(t, uuid)
}
