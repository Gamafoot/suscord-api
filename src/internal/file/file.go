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

type fileManager struct {
	config *config.Config
}

func NewFileManager(config *config.Config) *fileManager {
	return &fileManager{
		config: config,
	}
}

func (m *fileManager) Save(filename string, content []byte) (string, error) {
	ext := filepath.Ext(filename)
	filename = fmt.Sprintf("%d"+ext, time.Now().UnixNano())

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

	rootpath = fmt.Sprintf("%s/%d/%d", m.config.Media.RootFolder, year, month)
	filepath := fmt.Sprintf("%s/%s", rootpath, filename)

	os.MkdirAll(rootpath, os.ModePerm)

	err := os.WriteFile(filepath, content, os.ModePerm)
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}

	return filepath, nil
}

func (m *fileManager) SaveFileHeader(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", pkgErrors.WithStack(err)
	}

	return m.Save(fileHeader.Filename, content)
}
