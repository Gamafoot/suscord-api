package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
	"suscord/pkg/urlpath"
)

func NewAttachment(attachment *entity.Attachment, mediaURL string) *dto.Attachment {
	return &dto.Attachment{
		ID:       attachment.ID,
		FileUrl:  urlpath.GetMediaURL(mediaURL, attachment.FilePath),
		FileSize: attachment.FileSize,
		MimeType: attachment.MimeType,
	}
}
