package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"suscord/internal/config"
	brokermocks "suscord/internal/domain/broker/mocks"
	brokermsg "suscord/internal/domain/broker/message"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	cachemocks "suscord/internal/domain/storage/cache/mocks"
	dbmocks "suscord/internal/domain/storage/database/mocks"
	storagemocks "suscord/internal/domain/storage/mocks"
)

func TestChatMemberService_IsMemberOfChat_Delegates(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	repo := dbmocks.NewMockChatMemberStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(repo).Maybe()

	repo.EXPECT().IsMemberOfChat(mock.Anything, userID, chatID).Return(true, nil).Once()

	s := NewChatMemberService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	ok, err := s.IsMemberOfChat(ctx, userID, chatID)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestChatMemberService_GetNonMembers_OK(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	memberID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	repo := dbmocks.NewMockChatMemberStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(repo).Maybe()

	repo.EXPECT().IsMemberOfChat(mock.Anything, memberID, chatID).Return(true, nil).Once()
	repo.EXPECT().GetChatMembers(mock.Anything, chatID).Return([]*entity.User{{ID: 1}}, nil).Once()

	s := NewChatMemberService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	users, err := s.GetNonMembers(ctx, chatID, memberID)
	require.NoError(t, err)
	require.Len(t, users, 1)
}

func TestChatMemberService_GetNonMembers_NotMember(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	memberID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	repo := dbmocks.NewMockChatMemberStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(repo).Maybe()

	repo.EXPECT().IsMemberOfChat(mock.Anything, memberID, chatID).Return(false, nil).Once()

	s := NewChatMemberService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	users, err := s.GetNonMembers(ctx, chatID, memberID)
	assert.Nil(t, users)
	assert.ErrorIs(t, err, domainErrors.ErrUserIsNotMemberOfChat)
}

func TestChatMemberService_SendInvite_OK_PutsInCacheAndPublishes(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	ownerID := uint(1)
	chatID := uint(1)
	inviteeID := uint(2)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)
	cache := cachemocks.NewMockStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	storage.EXPECT().Cache().Return(cache).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatMemberRepo.EXPECT().IsMemberOfChat(mock.Anything, ownerID, chatID).Return(true, nil).Once()
	chatRepo.EXPECT().GetByID(mock.Anything, chatID).Return(&entity.Chat{ID: chatID, Type: "group"}, nil).Once()

	var code string
	cache.EXPECT().
		Set(mock.Anything, mock.AnythingOfType("string"), chatID, mock.AnythingOfType("time.Duration")).
		Run(func(_ context.Context, key string, value any, ttl time.Duration) {
			code = key
			assert.Equal(t, chatID, value)
			assert.Equal(t, 10*time.Second, ttl)
		}).
		Return(nil).
		Once()

	broker.EXPECT().
		Publish(mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ brokermsg.BrokerMessage) { assert.NotEmpty(t, code) }).
		Return(nil).
		Once()

	s := NewChatMemberService(cfg, storage, broker, log)
	err := s.SendInvite(ctx, ownerID, chatID, inviteeID)
	require.NoError(t, err)
}

func TestChatMemberService_SendInvite_ChatNotGroup(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	ownerID := uint(1)
	chatID := uint(1)
	inviteeID := uint(2)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatMemberRepo.EXPECT().IsMemberOfChat(mock.Anything, ownerID, chatID).Return(true, nil).Once()
	chatRepo.EXPECT().GetByID(mock.Anything, chatID).Return(&entity.Chat{ID: chatID, Type: "private"}, nil).Once()

	s := NewChatMemberService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	err := s.SendInvite(ctx, ownerID, chatID, inviteeID)
	assert.ErrorIs(t, err, domainErrors.ErrChatIsNotGroup)
}

func TestChatMemberService_AcceptInvite_OK(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Media.Url = "/media/"

	userID := uint(1)
	chatID := uint(1)
	code := "code"
	chatIDStr := "1"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)
	cache := cachemocks.NewMockStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	storage.EXPECT().Cache().Return(cache).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	cache.EXPECT().Get(mock.Anything, code).Return(chatIDStr, nil).Once()
	chatMemberRepo.EXPECT().AddUserToChat(mock.Anything, userID, chatID).Return(nil).Once()
	userRepo.EXPECT().GetByID(mock.Anything, userID).Return(&entity.User{ID: userID, Username: "u"}, nil).Once()

	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()

	s := NewChatMemberService(cfg, storage, broker, log)
	err := s.AcceptInvite(ctx, userID, code)
	require.NoError(t, err)
}

func TestChatMemberService_AcceptInvite_InvalidChatIDInCache(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	code := "code"

	storage := storagemocks.NewMockStorage(t)
	cache := cachemocks.NewMockStorage(t)

	storage.EXPECT().Cache().Return(cache).Maybe()
	cache.EXPECT().Get(mock.Anything, code).Return("not-a-number", nil).Once()

	s := NewChatMemberService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	err := s.AcceptInvite(ctx, userID, code)
	require.Error(t, err)
}

func TestChatMemberService_LeaveFromChat_PublishErrorIsLogged(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()

	chatMemberRepo.EXPECT().IsMemberOfChat(mock.Anything, userID, chatID).Return(true, nil).Once()
	chatMemberRepo.EXPECT().Delete(mock.Anything, userID, chatID).Return(nil).Once()

	pubErr := errors.New("broker down")
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(pubErr).Once()
	log.On("Err", pubErr, mock.AnythingOfType("[]logger.Field")).Return().Once()

	s := NewChatMemberService(cfg, storage, broker, log)
	err := s.LeaveFromChat(ctx, chatID, userID)
	require.NoError(t, err)
}
