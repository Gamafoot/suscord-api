package storage

import (
	"suscord/internal/database/gorm/model"
	"suscord/internal/domain/entity"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type messageAttachmentStorage struct {
	db *gorm.DB
}

func NewMessageAttachmentStorage(db *gorm.DB) *messageAttachmentStorage {
	return &messageAttachmentStorage{db: db}
}

func (repo *messageAttachmentStorage) GetByID(attachmentID int) (*entity.MessageAttachment, error) {
	attachment := new(model.MessageAttachment)
	if err := repo.db.First(&attachment, "id = ?", attachmentID).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return messageAttachmentModelToDomain(attachment), nil
}

func (repo *messageAttachmentStorage) GetByMessageID(messageID int) (*entity.MessageAttachment, error) {
	attachment := new(model.MessageAttachment)
	if err := repo.db.First(&attachment, "message_id = ?", messageID).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return messageAttachmentModelToDomain(attachment), nil
}

func (repo *messageAttachmentStorage) Create(attachment *entity.MessageAttachment) error {
	attachmentModel := messageAttachmentDomainToModel(attachment)
	if err := repo.db.Create(attachmentModel).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (repo *messageAttachmentStorage) Update(attachmentID uint, data map[string]interface{}) error {
	err := repo.db.Model(&model.MessageAttachment{}).Where("id = ?", attachmentID).Updates(data).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (repo *messageAttachmentStorage) Delete(attachmentID uint) error {
	if err := repo.db.Delete(&entity.MessageAttachment{ID: attachmentID}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func messageAttachmentModelToDomain(attachment *model.MessageAttachment) *entity.MessageAttachment {
	return &entity.MessageAttachment{
		ID:       attachment.ID,
		FileURL:  attachment.FileURL,
		FileSize: attachment.FileSize,
		MimeType: attachment.MimeType,
	}
}

func messageAttachmentDomainToModel(attachment *entity.MessageAttachment) *model.MessageAttachment {
	return &model.MessageAttachment{
		ID:       attachment.ID,
		FileURL:  attachment.FileURL,
		FileSize: attachment.FileSize,
		MimeType: attachment.MimeType,
	}
}
