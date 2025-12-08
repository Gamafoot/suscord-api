package storage

import (
	"suscord/internal/domain/storage/cache"
	"suscord/internal/domain/storage/database"
	"suscord/internal/domain/storage/file"
)

type Storage interface {
	Database() database.Storage
	Cache() cache.Storage
	File() file.FileStorage
}
