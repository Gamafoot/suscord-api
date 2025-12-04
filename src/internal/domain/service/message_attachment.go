package service

import (
	"context"
)

type AttachementService interface {
	Delete(ctx context.Context, userID, attachmentID uint) error
}
