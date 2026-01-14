package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"suscord/internal/config"
	brokermocks "suscord/internal/domain/broker/mocks"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	dbmocks "suscord/internal/domain/storage/database/mocks"
	filemocks "suscord/internal/domain/storage/file/mocks"
	storagemocks "suscord/internal/domain/storage/mocks"
)

func TestMessageService_GetChatMessages_OK(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)
	lastMessageID := uint(0)
	limit := 50
	messageID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	chatMemberRepo.EXPECT().
		IsMemberOfChat(mock.Anything, userID, chatID).
		Return(true, nil).
		Once()

	messageRepo.EXPECT().
		GetMessages(mock.Anything, chatID, lastMessageID, limit).
		Return([]*entity.Message{{ID: messageID}}, nil).
		Once()

	s := NewMessageService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	msgs, err := s.GetChatMessages(ctx, &entity.GetMessagesInput{UserID: userID, ChatID: chatID, LastMessageID: lastMessageID, Limit: limit})
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	assert.Equal(t, messageID, msgs[0].ID)
}

func TestMessageService_GetChatMessages_NotMember(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	chatMemberRepo := dbmocks.NewMockChatMemberStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().ChatMember().Return(chatMemberRepo).Maybe()

	chatMemberRepo.EXPECT().
		IsMemberOfChat(mock.Anything, userID, chatID).
		Return(false, nil).
		Once()

	s := NewMessageService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	msgs, err := s.GetChatMessages(ctx, &entity.GetMessagesInput{UserID: userID, ChatID: chatID})
	assert.Nil(t, msgs)
	assert.ErrorIs(t, err, domainErrors.ErrUserIsNotMemberOfChat)
}

func TestMessageService_Create_NoFiles_Publishes(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Media.Url = "/media/"

	userID := uint(1)
	chatID := uint(1)
	messageID := uint(1)
	createInput := &entity.CreateMessage{Type: "text", Content: "hi"}

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	messageRepo.EXPECT().
		Create(mock.Anything, userID, chatID, mock.AnythingOfType("*entity.CreateMessage")).
		Return(messageID, nil).
		Once()

	messageRepo.EXPECT().
		GetByID(mock.Anything, messageID).
		Return(&entity.Message{ID: messageID, ChatID: chatID}, nil).
		Once()

	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()

	s := NewMessageService(cfg, storage, broker, log)
	msg, err := s.Create(ctx, userID, chatID, createInput, nil)
	require.NoError(t, err)
	require.NotNil(t, msg)
	assert.Equal(t, messageID, msg.ID)
	assert.Empty(t, msg.Attachments)
}

func TestMessageService_Create_PublishError_IsLogged(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Media.Url = "/media/"

	userID := uint(1)
	chatID := uint(1)
	messageID := uint(1)
	createInput := &entity.CreateMessage{Type: "text", Content: "hi"}

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	messageRepo.EXPECT().
		Create(mock.Anything, userID, chatID, mock.Anything).
		Return(messageID, nil).
		Once()

	messageRepo.EXPECT().
		GetByID(mock.Anything, messageID).
		Return(&entity.Message{ID: messageID, ChatID: chatID}, nil).
		Once()

	pubErr := errors.New("broker down")
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(pubErr).Once()
	log.On("Err", pubErr, mock.AnythingOfType("[]logger.Field")).Return().Once()

	s := NewMessageService(cfg, storage, broker, log)
	msg, err := s.Create(ctx, userID, chatID, createInput, nil)
	require.NoError(t, err)
	require.NotNil(t, msg)
}

func TestMessageService_Create_WithFiles_AttachesUploadedFiles(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Media.Url = "/media/"

	userID := uint(1)
	chatID := uint(1)
	messageID := uint(1)
	attachmentID := uint(1)
	createInput := &entity.CreateMessage{Type: "text", Content: "hi"}

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)
	attachmentRepo := dbmocks.NewMockAttachmentStorage(t)
	fileStorage := filemocks.NewMockFileStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	storage.EXPECT().File().Return(fileStorage).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()
	db.EXPECT().Attachment().Return(attachmentRepo).Maybe()

	messageRepo.EXPECT().
		Create(mock.Anything, userID, chatID, mock.Anything).
		Return(messageID, nil).
		Once()
	messageRepo.EXPECT().
		GetByID(mock.Anything, messageID).
		Return(&entity.Message{ID: messageID, ChatID: chatID}, nil).
		Once()

	fh := &multipart.FileHeader{Filename: "a.png", Size: 123}
	fileStorage.EXPECT().UploadFile(fh, "messages").Return("/media/messages/a.png", nil).Once()
	attachmentRepo.EXPECT().
		Create(mock.Anything, messageID, mock.AnythingOfType("*entity.CreateAttachment")).
		Return(attachmentID, nil).
		Once()
	attachmentRepo.EXPECT().
		GetByID(mock.Anything, attachmentID).
		Return(&entity.Attachment{ID: attachmentID, FilePath: "/media/messages/a.png"}, nil).
		Once()

	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()

	s := NewMessageService(cfg, storage, broker, log)
	msg, err := s.Create(ctx, userID, chatID, createInput, []*multipart.FileHeader{fh})
	require.NoError(t, err)
	require.NotNil(t, msg)
	require.Len(t, msg.Attachments, 1)
	assert.Equal(t, attachmentID, msg.Attachments[0].ID)
}

func TestMessageService_Update_NotOwner(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	messageID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	messageRepo.EXPECT().
		IsOwner(mock.Anything, userID, messageID).
		Return(false, nil).
		Once()

	s := NewMessageService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	msg, err := s.Update(ctx, userID, messageID, &entity.UpdateMessage{Content: "x"})
	assert.Nil(t, msg)
	assert.ErrorIs(t, err, domainErrors.ErrUserIsNotMemberOfChat)
}

func TestMessageService_Delete_PublishError_IsLoggedButNotReturned(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Media.Url = "/media/"

	userID := uint(1)
	messageID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	messageRepo.EXPECT().IsOwner(mock.Anything, userID, messageID).Return(true, nil).Once()
	messageRepo.EXPECT().GetByID(mock.Anything, messageID).Return(&entity.Message{ID: messageID, ChatID: chatID}, nil).Once()
	messageRepo.EXPECT().Delete(mock.Anything, messageID).Return(nil).Once()

	pubErr := errors.New("broker down")
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(pubErr).Once()
	log.On("Err", pubErr, mock.AnythingOfType("[]logger.Field")).Return().Once()

	s := NewMessageService(cfg, storage, broker, log)
	err := s.Delete(ctx, userID, messageID)
	require.NoError(t, err)
}
