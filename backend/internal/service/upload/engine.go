package upload

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

func (e *engine) Upload(ctx context.Context, req *request.UploadRequest) (*vo.UploadResult, error) {
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
	return &vo.UploadResult{
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

func (e *engine) DeleteFile(ctx context.Context, storedPath string) error {
	return e.storage.Delete(storedPath)
}
