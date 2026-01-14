package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"suscord/internal/domain/entity"
	dbmocks "suscord/internal/domain/storage/database/mocks"
	storagemocks "suscord/internal/domain/storage/mocks"
)

func TestUserService_GetByID_OK(t *testing.T) {
	ctx := context.Background()

	userID := uint(1)
	username := "u"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	userRepo.EXPECT().
		GetByID(mock.Anything, userID).
		Return(&entity.User{ID: userID, Username: username}, nil).
		Once()

	s := NewUserService(storage)
	u, err := s.GetByID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, u.ID)
	assert.Equal(t, username, u.Username)
}

func TestUserService_GetByID_ErrorWrapped(t *testing.T) {
	ctx := context.Background()

	userID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	dbErr := errors.New("db down")
	userRepo.EXPECT().
		GetByID(mock.Anything, userID).
		Return(nil, dbErr).
		Once()

	s := NewUserService(storage)
	u, err := s.GetByID(ctx, userID)
	assert.Nil(t, u)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user")
	assert.ErrorIs(t, err, dbErr)
}

func TestUserService_SearchUsers_Delegates(t *testing.T) {
	ctx := context.Background()

	userID := uint(1)
	foundUserID := uint(2)
	search := "a"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	userRepo.EXPECT().
		SearchUsers(mock.Anything, userID, search).
		Return([]*entity.User{{ID: foundUserID, Username: "alex"}}, nil).
		Once()

	s := NewUserService(storage)
	users, err := s.SearchUsers(ctx, userID, search)
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, "alex", users[0].Username)
}
