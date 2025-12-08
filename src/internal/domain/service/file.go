package service

import "mime/multipart"

type FileService interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}
