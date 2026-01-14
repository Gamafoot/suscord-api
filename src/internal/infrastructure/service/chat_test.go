package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"suscord/internal/config"
	brokermocks "suscord/internal/domain/broker/mocks"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	dbmocks "suscord/internal/domain/storage/database/mocks"
	storagemocks "suscord/internal/domain/storage/mocks"
)

func TestChatService_GetUserChats_WithoutSearch(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatRepo.EXPECT().GetUserChats(mock.Anything, userID).Return([]*entity.Chat{{ID: chatID}}, nil).Once()

	s := NewChatService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	chats, err := s.GetUserChats(ctx, userID, "")
	require.NoError(t, err)
	require.Len(t, chats, 1)
}

func TestChatService_GetUserChats_WithSearch(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)
	search := "a"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatRepo.EXPECT().SearchUserChats(mock.Anything, userID, search).Return([]*entity.Chat{{ID: chatID}}, nil).Once()

	s := NewChatService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	chats, err := s.GetUserChats(ctx, userID, search)
	require.NoError(t, err)
	require.Len(t, chats, 1)
}

func TestChatService_GetOrCreatePrivateChat_ExistingChat(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Static.URL = "/static/"

	userID := uint(1)
	friendID := uint(2)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	broker := brokermocks.NewMockBroker(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatMemberRepo.EXPECT().GetPrivateChatID(mock.Anything, userID, friendID).Return(chatID, nil).Once()
	chatRepo.EXPECT().GetUserChat(mock.Anything, chatID, userID).Return(&entity.Chat{ID: chatID, Type: "private"}, nil).Once()

	s := NewChatService(cfg, storage, broker, &mockLogger{})
	chat, err := s.GetOrCreatePrivateChat(ctx, &entity.CreatePrivateChat{UserID: userID, FriendID: friendID})
	require.NoError(t, err)
	require.NotNil(t, chat)
	assert.Equal(t, chatID, chat.ID)
}

func TestChatService_GetOrCreatePrivateChat_NotFound_CreatesAndPublishes(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Static.URL = "/static/"

	userID := uint(1)
	friendID := uint(2)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)
	userRepo := dbmocks.NewMockUserStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()
	db.EXPECT().User().Return(userRepo).Maybe()

	chatMemberRepo.EXPECT().GetPrivateChatID(mock.Anything, userID, friendID).Return(uint(0), domainErrors.ErrRecordNotFound).Once()

	chatRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*entity.CreateChat")).Return(chatID, nil).Once()
	chatMemberRepo.EXPECT().AddUserToChat(mock.Anything, userID, chatID).Return(nil).Once()
	chatMemberRepo.EXPECT().AddUserToChat(mock.Anything, friendID, chatID).Return(nil).Once()

	chatRepo.EXPECT().GetUserChat(mock.Anything, chatID, userID).Return(&entity.Chat{ID: chatID, Type: "private"}, nil).Once()

	userRepo.EXPECT().GetByID(mock.Anything, userID).Return(&entity.User{ID: userID}, nil).Once()
	userRepo.EXPECT().GetByID(mock.Anything, friendID).Return(&entity.User{ID: friendID}, nil).Once()

	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Twice()

	s := NewChatService(cfg, storage, broker, log)
	chat, err := s.GetOrCreatePrivateChat(ctx, &entity.CreatePrivateChat{UserID: userID, FriendID: friendID})
	require.NoError(t, err)
	require.NotNil(t, chat)
	assert.Equal(t, chatID, chat.ID)
}

func TestChatService_UpdateGroupChat_PublishErrorIsLogged(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Media.Url = "/media/"

	userID := uint(1)
	chatID := uint(1)
	newName := "x"

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatRepo.EXPECT().Update(mock.Anything, chatID, mock.AnythingOfType("*entity.UpdateChat")).Return(nil).Once()
	chatRepo.EXPECT().GetByID(mock.Anything, chatID).Return(&entity.Chat{ID: chatID, Type: "group"}, nil).Once()

	pubErr := errors.New("broker down")
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(pubErr).Once()
	log.On("Err", pubErr, mock.AnythingOfType("[]logger.Field")).Return().Once()

	s := NewChatService(cfg, storage, broker, log)
	chat, err := s.UpdateGroupChat(ctx, userID, chatID, &entity.UpdateChat{Name: &newName})
	require.NoError(t, err)
	require.NotNil(t, chat)
}

func TestChatService_DeletePrivateChat_OK(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatMemberRepo.EXPECT().IsMemberOfChat(mock.Anything, userID, chatID).Return(true, nil).Once()
	chatRepo.EXPECT().GetByID(mock.Anything, chatID).Return(&entity.Chat{ID: chatID, Type: "private"}, nil).Once()
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()
	chatRepo.EXPECT().Delete(mock.Anything, chatID).Return(nil).Once()

	s := NewChatService(cfg, storage, broker, log)
	err := s.DeletePrivateChat(ctx, userID, chatID)
	require.NoError(t, err)
}

func TestChatService_DeletePrivateChat_NotMember(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()

	chatMemberRepo.EXPECT().IsMemberOfChat(mock.Anything, userID, chatID).Return(false, nil).Once()

	s := NewChatService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	err := s.DeletePrivateChat(ctx, userID, chatID)
	assert.ErrorIs(t, err, domainErrors.ErrForbidden)
}

func TestChatService_DeletePrivateChat_NotPrivate(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	chatRepo := dbmocks.NewMockChatStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Chat().Return(chatRepo).Maybe()

	chatMemberRepo.EXPECT().IsMemberOfChat(mock.Anything, userID, chatID).Return(true, nil).Once()
	chatRepo.EXPECT().GetByID(mock.Anything, chatID).Return(&entity.Chat{ID: chatID, Type: "group"}, nil).Once()

	s := NewChatService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	err := s.DeletePrivateChat(ctx, userID, chatID)
	assert.ErrorIs(t, err, domainErrors.ErrForbidden)
}
