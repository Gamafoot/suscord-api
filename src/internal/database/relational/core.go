package relational

import (
	implStorage "suscord/internal/database/relational/storage"
	"suscord/internal/domain/storage"

	errors "github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewConnect(dbURL, log_level string) (*gorm.DB, error) {
	var logger gormLogger.Interface

	switch log_level {
	case "info":
		logger = gormLogger.Default.LogMode(gormLogger.Info)
	case "warn":
		logger = gormLogger.Default.LogMode(gormLogger.Warn)
	case "error":
		logger = gormLogger.Default.LogMode(gormLogger.Error)
	default:
		logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewGormStorage(dbURL, log_level string) (storage.Storage, error) {
	db, err := NewConnect(dbURL, log_level)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return implStorage.NewGormStorage(db), nil
}
