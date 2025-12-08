package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"suscord/internal/config"
	"time"

	pkgErrors "github.com/pkg/errors"
)

type fileStorage struct {
	cfg *config.Config
}

func NewStorage(cfg *config.Config) *fileStorage {
	return &fileStorage{
		cfg: cfg,
	}
}

func (s *fileStorage) UploadFile(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d"+ext, file.Filename, time.Now().UnixNano())

	var (
		rootpath string
		year     int
		month    int
	)

	if year == 0 && month == 0 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	rootpath = fmt.Sprintf("%s/%d/%d", s.cfg.Media.RootFolder, year, month)
	filepath := fmt.Sprintf("%s/%s", rootpath, filename)

	os.MkdirAll(rootpath, os.ModePerm)

	src, err := file.Open()
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}
	defer dst.Close()

	io.Copy(dst, src)

	return filepath, nil
}
