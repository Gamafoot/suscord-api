package app

import (
	"suscord/internal/config"
	"suscord/internal/domain/storage"
	"suscord/internal/domain/storage/cache"
	"suscord/internal/domain/storage/database"
	"suscord/internal/domain/storage/file"
	"suscord/internal/infrastructure/cache/redis"
	"suscord/internal/infrastructure/database/relational"
	fileStorage "suscord/internal/infrastructure/file"

	pkgErrors "github.com/pkg/errors"
)

type storageImpl struct {
	dbStorage database.Storage
	cache     cache.Storage
	file      file.FileStorage
}

func NewStorage(cfg *config.Config) (storage.Storage, error) {
	dbStorage, err := relational.NewStorage(cfg.Database.URL, cfg.Database.LogLevel)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	redisStorage := redis.NewStorage(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	fileStorage := fileStorage.NewStorage(cfg)

	return &storageImpl{
		dbStorage: dbStorage,
		cache:     redisStorage,
		file:      fileStorage,
	}, nil
}

func (s *storageImpl) Database() database.Storage {
	return s.dbStorage
}

func (s *storageImpl) Cache() cache.Storage {
	return s.cache
}

func (s *storageImpl) File() file.FileStorage {
	return s.file
}
