// package upload

// import (
// 	"context"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"path/filepath"
// 	"strings"
// 	"time"
// 	"tiny-forum/internal/model/do"

// 	"github.com/google/uuid"
// )

// func (s *engine) Upload(ctx context.Context, req *UploadRequest) (*do.Attachment, error) {
// 	// 1. 打开文件
// 	src, err := req.File.Open()
// 	if err != nil {
// 		return nil, fmt.Errorf("open file: %w", err)
// 	}
// 	defer src.Close()

// 	// 2. 计算哈希
// 	hash, err := s.computeHash(src)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fileHash := hex.EncodeToString(hash)
// 	// 重置reader
// 	src.Seek(0, io.SeekStart)

// 	// 3. 去重检查
// 	dup, err := s.attachmentRepo.FindDuplicate(ctx, fileHash, req.FileType)
// 	if err == nil {
// 		// 已存在相同文件，直接返回（不增加引用计数，简单场景）
// 		return dup, nil
// 	}

// 	// 4. 获取对应策略
// 	handler, err := s.registry.Get(req.FileType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 5. 校验文件
// 	if err := handler.Validate(req.File); err != nil {
// 		return nil, err
// 	}

// 	// 6. 构建元数据
// 	meta := &do.Attachment{
// 		FileID:       uuid.New().String(),
// 		UserID:       req.UserID,
// 		PluginID:     req.PluginID,
// 		PostID:       req.PostID,
// 		ReplyID:      req.ReplyID,
// 		OriginalName: req.File.Filename,
// 		StoredName:   s.generateStoredName(req.File.Filename),
// 		Size:         req.File.Size,
// 		FileType:     req.FileType,
// 		MimeType:     req.File.Header.Get("Content-Type"),
// 		MimeMajor:    s.extractMimeMajor(req.File.Header.Get("Content-Type")),
// 		Ext:          strings.TrimPrefix(filepath.Ext(req.File.Filename), "."),
// 		Status:       do.StatusTemp,
// 		UploadIP:     req.ClientIP,
// 		FileHash:     fileHash,
// 	}

// 	// 7. 处理文件（可选转换）
// 	processedReader, err := handler.Process(src, meta)
// 	if err != nil {
// 		return nil, fmt.Errorf("process file: %w", err)
// 	}

// 	// 8. 存储路径
// 	storagePath := handler.GetStoragePath(meta)

// 	// 9. 保存到驱动
// 	if _, err := s.storage.Save(processedReader, storagePath); err != nil {
// 		return nil, fmt.Errorf("save file: %w", err)
// 	}
// 	meta.StoredPath = storagePath

// 	// 10. 保存元数据到数据库
// 	if err := s.attachmentRepo.Create(ctx, meta); err != nil {
// 		// 回滚已存储文件
// 		_ = s.storage.Delete(storagePath)
// 		return nil, fmt.Errorf("create record: %w", err)
// 	}

// 	// 11. 异步更新状态为正常（也可在后台任务中处理）
// 	go func() {
// 		// 模拟可能的后处理（如缩略图生成），完成后更新状态
// 		time.Sleep(200 * time.Millisecond)
// 		meta.Status = do.StatusNormal
// 		_ = s.attachmentRepo.Update(context.Background(), meta)
// 	}()

// 	return meta, nil
// }

// func (s *engine) computeHash(r io.Reader) ([]byte, error) {
// 	h := sha256.New()
// 	if _, err := io.Copy(h, r); err != nil {
// 		return nil, err
// 	}
// 	return h.Sum(nil), nil
// }

// func (s *engine) generateStoredName(original string) string {
// 	ext := strings.ToLower(filepath.Ext(original))
// 	return uuid.New().String() + ext
// }

//	func (s *engine) extractMimeMajor(mime string) do.MimeTypeMajor {
//		switch {
//		case strings.HasPrefix(mime, "image/"):
//			return do.MimeImage
//		case strings.HasPrefix(mime, "video/"):
//			return do.MimeVideo
//		case strings.HasPrefix(mime, "audio/"):
//			return do.MimeAudio
//		case strings.Contains(mime, "pdf") || strings.Contains(mime, "document") || strings.Contains(mime, "text"):
//			return do.MimeDocument
//		default:
//			return do.MimeOther
//		}
//	}
package upload