package file

import "mime/multipart"

type FileStorage interface {
	UploadFile(file *multipart.FileHeader, uploadTo ...string) (string, error)
}
