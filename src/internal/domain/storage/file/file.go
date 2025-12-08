package file

import "mime/multipart"

type FileStorage interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}
