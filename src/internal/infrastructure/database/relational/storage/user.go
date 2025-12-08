package storage

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/infrastructure/database/relational/model"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *userStorage {
	return &userStorage{db: db}
}

func (s *userStorage) GetByID(ctx context.Context, userID uint) (*entity.User, error) {
	user := new(model.User)
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}
	return userModelToEntity(user), nil
}

func (s *userStorage) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(model.User)
	if err := s.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}
	return userModelToEntity(user), nil
}

func (s *userStorage) SearchUsers(ctx context.Context, exceptUserID uint, username string) ([]*entity.User, error) {
	users := make([]*model.User, 0)
	if err := s.db.WithContext(ctx).Order("username ASC").Find(&users, "id != ? AND username ~* ?", exceptUserID, username).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	result := make([]*entity.User, len(users))

	for i, user := range users {
		result[i] = userModelToEntity(user)
	}

	return result, nil
}

func (s *userStorage) Create(ctx context.Context, user *entity.User) error {
	userModel := userDomainToEntity(user)
	if err := s.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *userStorage) Update(ctx context.Context, userID uint, data map[string]interface{}) error {
	err := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(data).Error
	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return domainErrors.ErrRecordNotFound
		}
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *userStorage) Delete(ctx context.Context, userID uint) error {
	if err := s.db.WithContext(ctx).Delete(&entity.User{ID: userID}).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return domainErrors.ErrRecordNotFound
		}
		return pkgErrors.WithStack(err)
	}
	return nil
}

func userModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		ID:         user.ID,
		Username:   user.Username,
		Password:   user.Password,
		AvatarPath: user.AvatarPath,
		FriendCode: user.FriendCode,
	}
}

func userDomainToEntity(user *entity.User) *model.User {
	return &model.User{
		ID:         user.ID,
		Username:   user.Username,
		Password:   user.Password,
		AvatarPath: user.AvatarPath,
		FriendCode: user.FriendCode,
	}
}
