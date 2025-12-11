package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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

func (s *fileStorage) UploadFile(file *multipart.FileHeader, uploadTo ...string) (string, error) {
	var uploadPath string

	if len(uploadTo) > 0 {
		uploadPath = uploadTo[0]
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	rootpath := filepath.Join(s.cfg.Media.Folder, uploadPath)
	filePath := filepath.Join(rootpath, filename)

	os.MkdirAll(rootpath, os.ModePerm)

	src, err := file.Open()
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}
	defer dst.Close()

	io.Copy(dst, src)

	filePath, _ = strings.CutPrefix(filePath, s.cfg.Media.Folder)

	return filePath, nil
}
