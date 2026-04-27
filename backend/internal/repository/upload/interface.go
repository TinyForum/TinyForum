package upload

import (
	"context"
	"gorm.io/gorm"
	"time"
	"tiny-forum/internal/model"
)

type UploadRepository interface {
	Create(ctx context.Context, attachment *model.Attachment) error
	GetByFileID(ctx context.Context, fileID string) (*model.Attachment, error)
	GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*model.Attachment, error)
	Delete(ctx context.Context, fileID string) error
	UpdateStatus(ctx context.Context, fileID string, status int) error
	ListByUser(ctx context.Context, userID int64, fileType string, limit, offset int) ([]*model.Attachment, int64, error)
	DeleteUnusedTemp(ctx context.Context, beforeTime time.Time) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) UploadRepository {
	return &repository{db: db}
}
