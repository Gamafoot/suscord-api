package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"suscord/internal/domain/logger"
	"time"

	pkgErrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	entry *logrus.Logger
	close func()
}

func NewLogger(level string, folderPath string) (*Logger, error) {
	if folderPath == "" {
		return nil, pkgErrors.New("folderPath пустой")
	}

	_, err := os.Stat(folderPath)
	if pkgErrors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(folderPath, 0755); err != nil {
			return nil, pkgErrors.WithStack(err)
		}
	}

	writer, close, err := newWriter(folderPath)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetOutput(writer)

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	Level, err := logrus.ParseLevel(level)
	if err != nil {
		Level = logrus.InfoLevel
	}
	logger.SetLevel(Level)

	return &Logger{
		entry: logger,
		close: close,
	}, nil
}

func (l *Logger) Debug(msg string, fields ...logger.Field) {
	l.entry.WithFields(toLogrusFields(fields)).Debug(msg)
}

func (l *Logger) Info(msg string, fields ...logger.Field) {
	l.entry.WithFields(toLogrusFields(fields)).Info(msg)
}

func (l *Logger) Warn(msg string, fields ...logger.Field) {
	l.entry.WithFields(toLogrusFields(fields)).Warn(msg)
}

func (l *Logger) Error(msg string, fields ...logger.Field) {
	l.entry.WithFields(toLogrusFields(fields)).Error(msg)
}

func (l *Logger) Err(err error, fields ...logger.Field) {
	l.entry.WithFields(toLogrusFields(fields)).Error(fmt.Sprintf("%+v\n", err))
}

func toLogrusFields(fields []logger.Field) logrus.Fields {
	f := logrus.Fields{}
	for _, field := range fields {
		f[field.Key] = field.Value
	}
	return f
}

func newWriter(logFolder string) (io.Writer, func(), error) {
	writers := []io.Writer{os.Stdout}

	var file *os.File
	var err error

	filepath := filepath.Join(logFolder, getLogFilename())

	if logFolder != "" {
		file, err = os.OpenFile(
			filepath,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, nil, pkgErrors.WithStack(err)
		}
		writers = append(writers, file)
	}

	close := func() {
		if file != nil {
			_ = file.Close()
		}
	}

	return io.MultiWriter(writers...), close, nil
}

func getLogFilename() string {
	t := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("log_%s.log", t.Format("2006_01_02"))
}
