package storage

import (
	"context"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/infrastructure/database/relational/model"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type attachmentStorage struct {
	db *gorm.DB
}

func NewAttachmentStorage(db *gorm.DB) *attachmentStorage {
	return &attachmentStorage{db: db}
}

func (s *attachmentStorage) GetByID(ctx context.Context, attachmentID uint) (*entity.Attachment, error) {
	attachment := new(model.Attachment)
	if err := s.db.WithContext(ctx).First(&attachment, "id = ?", attachmentID).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return attachmentModelToDomain(attachment), nil
}

func (s *attachmentStorage) GetByMessageID(ctx context.Context, messageID uint) ([]*entity.Attachment, error) {
	attachments := make([]*model.Attachment, 0)
	if err := s.db.WithContext(ctx).Find(&attachments, "message_id = ?", messageID).Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrRecordNotFound
		}
		return nil, pkgErrors.WithStack(err)
	}

	attachmentDomains := make([]*entity.Attachment, len(attachments))
	for i, attachment := range attachments {
		attachmentDomains[i] = attachmentModelToDomain(attachment)
	}

	return attachmentDomains, nil
}

func (s *attachmentStorage) Create(ctx context.Context, messageID uint, data *entity.CreateAttachment) (uint, error) {
	attachment := &model.Attachment{
		MessageID: messageID,
		FilePath:  data.FilePath,
		FileSize:  data.FileSize,
		MimeType:  data.MimeType,
	}
	if err := s.db.WithContext(ctx).Create(attachment).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}
	return attachment.ID, nil
}

func (s *attachmentStorage) Delete(ctx context.Context, attachmentID uint) error {
	if err := s.db.WithContext(ctx).Delete(&model.Attachment{ID: attachmentID}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *attachmentStorage) IsOwner(ctx context.Context, userID, attachmentID uint) (bool, error) {
	exists := false
	sql := `
	SELECT EXISTS(
		SELECT 1 FROM message_attachments 
		WHERE id = ? AND message_id IN (
			SELECT id FROM messages 
			WHERE user_id = ?
		)
	)`
	err := s.db.WithContext(ctx).Raw(sql, attachmentID, userID).Scan(&exists).Error
	if err != nil {
		return false, pkgErrors.WithStack(err)
	}
	return exists, nil
}

func attachmentModelToDomain(attachment *model.Attachment) *entity.Attachment {
	return &entity.Attachment{
		ID:       attachment.ID,
		FilePath: attachment.FilePath,
		FileSize: attachment.FileSize,
		MimeType: attachment.MimeType,
	}
}
