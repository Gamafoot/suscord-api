package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"suscord/internal/config"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/domain/storage"
	"time"

	pkgErrors "github.com/pkg/errors"
)

type attachmentService struct {
	cfg     *config.Config
	storage storage.Storage
}

func NewAttachmentService(cfg *config.Config, storage storage.Storage) *attachmentService {
	return &attachmentService{
		cfg:     cfg,
		storage: storage,
	}
}

func (s *attachmentService) Delete(ctx context.Context, userID, attachmentID uint) error {
	ok, err := s.storage.Attachment().IsOwner(ctx, userID, attachmentID)
	if err != nil {
		return err
	}

	if !ok {
		return domainErrors.ErrUserIsNotMemberOfChat
	}

	err = s.storage.Attachment().Delete(ctx, attachmentID)
	if err != nil {
		return err
	}

	return nil
}

func (s *attachmentService) saveFile(file *multipart.FileHeader) (string, error) {
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
