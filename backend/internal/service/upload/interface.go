package upload

import (
	"context"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	"tiny-forum/internal/storage"
	"tiny-forum/internal/strategy"
)

type Engine interface {
	// Upload 仅负责存储文件，返回存储结果，不操作数据库
	Upload(ctx context.Context, req *request.UploadRequest) (*vo.UploadResult, error)
	DeleteFile(ctx context.Context, storedPath string) error
}
type engine struct {
	storage  storage.StorageDriver
	registry *strategy.HandlerRegistry
}

func NewEngine(storage storage.StorageDriver, registry *strategy.HandlerRegistry) Engine {
	return &engine{
		storage:  storage,
		registry: registry,
	}
}
