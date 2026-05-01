package upload

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type UploadRepository interface {
	Create(ctx context.Context, attachment *do.Attachment) error
	GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error)
	GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*do.Attachment, error)
	Delete(ctx context.Context, fileID string) error
	UpdateStatus(ctx context.Context, fileID string, status int) error
	ListByUser(ctx context.Context, userID int64, fileType string, limit, offset int) ([]*do.Attachment, int64, error)
	DeleteUnusedTemp(ctx context.Context, beforeTime time.Time) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) UploadRepository {
	return &repository{db: db}
}
