package gorm

import (
	implStorage "suscord/internal/database/gorm/storage"
	"suscord/internal/domain/storage"

	errors "github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newConnect(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewGormStorage(dbURL string) (storage.Storage, error) {
	db, err := newConnect(dbURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return implStorage.NewGormStorage(db), nil
}
