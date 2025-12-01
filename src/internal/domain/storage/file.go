package storage

import (
	"mime/multipart"
)

type FileStorage interface {
	Save(filename string, file []byte) (string, error)
	SaveFileHeader(fileHeader *multipart.FileHeader) (string, error)
}
