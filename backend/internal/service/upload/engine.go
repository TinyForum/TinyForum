// internal/service/upload/engine.go
package upload

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/storage"
	"tiny-forum/internal/strategy"

	"github.com/google/uuid"
)

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

func (e *engine) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	src, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	hash, err := computeHash(src)
	if err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hash)
	src.Seek(0, io.SeekStart)

	handler, err := e.registry.Get(req.FileType)
	if err != nil {
		return nil, err
	}
	if err := handler.Validate(req.File); err != nil {
		return nil, err
	}

	// 构建临时元数据（仅用于存储路径生成）
	meta := &do.Attachment{
		UserID:       req.UserID,
		PluginID:     req.PluginID,
		PostID:       req.PostID,
		ReplyID:      req.ReplyID,
		OriginalName: req.File.Filename,
		StoredName:   generateStoredName(req.File.Filename),
		Size:         req.File.Size,
		FileType:     req.FileType,
		MimeType:     req.File.Header.Get("Content-Type"),
		MimeMajor:    extractMimeMajor(req.File.Header.Get("Content-Type")),
		Ext:          strings.TrimPrefix(filepath.Ext(req.File.Filename), "."),
		FileHash:     fileHash,
	}

	processed, err := handler.Process(src, meta)
	if err != nil {
		return nil, fmt.Errorf("process file: %w", err)
	}
	storagePath := handler.GetStoragePath(meta)
	if _, err := e.storage.Save(processed, storagePath); err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	// 返回结果，由调用方负责保存元数据到数据库
	return &UploadResult{
		FileHash:     fileHash,
		StoredPath:   storagePath,
		StoredName:   meta.StoredName,
		MimeType:     meta.MimeType,
		MimeMajor:    meta.MimeMajor,
		Ext:          meta.Ext,
		Size:         meta.Size,
		OriginalName: meta.OriginalName,
	}, nil
}

func computeHash(r io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func generateStoredName(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	return uuid.New().String() + ext
}

func extractMimeMajor(mime string) do.MimeTypeMajor {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return do.MimeImage
	case strings.HasPrefix(mime, "video/"):
		return do.MimeVideo
	case strings.HasPrefix(mime, "audio/"):
		return do.MimeAudio
	case strings.Contains(mime, "pdf") || strings.Contains(mime, "document") || strings.Contains(mime, "text"):
		return do.MimeDocument
	default:
		return do.MimeOther
	}
}
func (e *engine) DeleteFile(ctx context.Context, storedPath string) error {
	return e.storage.Delete(storedPath)
}