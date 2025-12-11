package service

import (
	"mime/multipart"
	"suscord/internal/domain/storage"
)

type fileService struct {
	storage storage.Storage
}

func NewFileService(storage storage.Storage) *fileService {
	return &fileService{
		storage: storage,
	}
}

func (s *fileService) UploadFile(file *multipart.FileHeader, uploadTo ...string) (string, error) {
	return s.storage.File().UploadFile(file, uploadTo...)
}
