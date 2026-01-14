package service

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	filemocks "suscord/internal/domain/storage/file/mocks"
	storagemocks "suscord/internal/domain/storage/mocks"
)

func TestFileService_UploadFile_NoUploadTo(t *testing.T) {
	filename := "a.png"
	uploadedPath := "/media/a.png"
	size := int64(1)

	storage := storagemocks.NewMockStorage(t)
	fileStorage := filemocks.NewMockFileStorage(t)

	storage.EXPECT().File().Return(fileStorage).Maybe()

	fh := &multipart.FileHeader{Filename: filename, Size: size}
	fileStorage.EXPECT().
		UploadFile(fh).
		Return(uploadedPath, nil).
		Once()

	s := NewFileService(storage)
	path, err := s.UploadFile(fh)
	require.NoError(t, err)
	assert.Equal(t, uploadedPath, path)
}

func TestFileService_UploadFile_WithUploadTo(t *testing.T) {
	filename := "a.png"
	uploadTo := "messages"
	uploadedPath := "/media/messages/a.png"
	size := int64(1)

	storage := storagemocks.NewMockStorage(t)
	fileStorage := filemocks.NewMockFileStorage(t)

	storage.EXPECT().File().Return(fileStorage).Maybe()

	fh := &multipart.FileHeader{Filename: filename, Size: size}
	fileStorage.EXPECT().
		UploadFile(fh, uploadTo).
		Return(uploadedPath, nil).
		Once()

	s := NewFileService(storage)
	path, err := s.UploadFile(fh, uploadTo)
	require.NoError(t, err)
	assert.Equal(t, uploadedPath, path)
}

func TestFileService_UploadFile_PropagatesError(t *testing.T) {
	filename := "a.png"
	uploadTo := "messages"
	size := int64(1)

	storage := storagemocks.NewMockStorage(t)
	fileStorage := filemocks.NewMockFileStorage(t)

	storage.EXPECT().File().Return(fileStorage).Maybe()

	fh := &multipart.FileHeader{Filename: filename, Size: size}
	uploadErr := errors.New("upload failed")
	fileStorage.EXPECT().
		UploadFile(fh, uploadTo).
		Return("", uploadErr).
		Once()

	s := NewFileService(storage)
	path, err := s.UploadFile(fh, uploadTo)
	assert.Empty(t, path)
	assert.ErrorIs(t, err, uploadErr)
}
