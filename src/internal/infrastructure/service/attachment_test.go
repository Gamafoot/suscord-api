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

func TestAttachmentService_Delete_OK(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	attachmentID := uint(1)
	messageID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	attachmentRepo := dbmocks.NewMockAttachmentStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Attachment().Return(attachmentRepo).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	attachmentRepo.EXPECT().IsOwner(mock.Anything, userID, attachmentID).Return(true, nil).Once()
	messageRepo.EXPECT().GetByAttachmentID(mock.Anything, attachmentID).Return(&entity.Message{ID: messageID, ChatID: chatID}, nil).Once()
	attachmentRepo.EXPECT().Delete(mock.Anything, attachmentID).Return(nil).Once()
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()

	s := NewAttachmentService(cfg, storage, broker, log)
	err := s.Delete(ctx, userID, attachmentID)
	require.NoError(t, err)
}

func TestAttachmentService_Delete_NotOwner(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	attachmentID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	attachmentRepo := dbmocks.NewMockAttachmentStorage(t)

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Attachment().Return(attachmentRepo).Maybe()

	attachmentRepo.EXPECT().IsOwner(mock.Anything, userID, attachmentID).Return(false, nil).Once()

	s := NewAttachmentService(cfg, storage, brokermocks.NewMockBroker(t), &mockLogger{})
	err := s.Delete(ctx, userID, attachmentID)
	assert.ErrorIs(t, err, domainErrors.ErrUserIsNotMemberOfChat)
}

func TestAttachmentService_Delete_PublishError_IsLoggedButNotReturned(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}

	userID := uint(1)
	attachmentID := uint(1)
	messageID := uint(1)
	chatID := uint(1)

	storage := storagemocks.NewMockStorage(t)
	db := dbmocks.NewMockStorage(t)
	attachmentRepo := dbmocks.NewMockAttachmentStorage(t)
	messageRepo := dbmocks.NewMockMessageStorage(t)

	broker := brokermocks.NewMockBroker(t)
	log := &mockLogger{}

	storage.EXPECT().Database().Return(db).Maybe()
	db.EXPECT().Attachment().Return(attachmentRepo).Maybe()
	db.EXPECT().Message().Return(messageRepo).Maybe()

	attachmentRepo.EXPECT().IsOwner(mock.Anything, userID, attachmentID).Return(true, nil).Once()
	messageRepo.EXPECT().GetByAttachmentID(mock.Anything, attachmentID).Return(&entity.Message{ID: messageID, ChatID: chatID}, nil).Once()
	attachmentRepo.EXPECT().Delete(mock.Anything, attachmentID).Return(nil).Once()

	pubErr := errors.New("broker down")
	broker.EXPECT().Publish(mock.Anything, mock.Anything).Return(pubErr).Once()
	log.On("Err", pubErr, mock.AnythingOfType("[]logger.Field")).Return().Once()

	s := NewAttachmentService(cfg, storage, broker, log)
	err := s.Delete(ctx, userID, attachmentID)
	require.NoError(t, err)
}
