package storage

import (
	"context"
	"suscord/internal/database/gorm/model"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *userStorage {
	return &userStorage{db: db}
}

func (repo *userStorage) GetByID(ctx context.Context, userID uint) (*entity.User, error) {
	user := new(model.User)
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}
	return userModelToDomain(user), nil
}

func (repo *userStorage) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(model.User)
	if err := repo.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}
	return userModelToDomain(user), nil
}

func (repo *userStorage) Create(ctx context.Context, user *entity.User) error {
	userModel := userDomainToModel(user)
	if err := repo.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (repo *userStorage) Update(ctx context.Context, userID uint, data map[string]interface{}) error {
	err := repo.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(data).Error
	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return domainErrors.ErrRecordNotFound
		}
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (repo *userStorage) Delete(ctx context.Context, userID uint) error {
	if err := repo.db.WithContext(ctx).Delete(&entity.User{ID: userID}).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return domainErrors.ErrRecordNotFound
		}
		return pkgErrors.WithStack(err)
	}
	return nil
}

func userModelToDomain(user *model.User) *entity.User {
	return &entity.User{
		ID:         user.ID,
		Username:   user.Username,
		Password:   user.Password,
		AvatarPath: user.AvatarPath,
		FriendCode: user.FriendCode,
	}
}

func userDomainToModel(user *entity.User) *model.User {
	return &model.User{
		ID:         user.ID,
		Username:   user.Username,
		Password:   user.Password,
		AvatarPath: user.AvatarPath,
		FriendCode: user.FriendCode,
	}
}
