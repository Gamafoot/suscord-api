package service

import "mime/multipart"

type FileService interface {
	UploadFile(file *multipart.FileHeader, fileUpload ...string) (string, error)
}
