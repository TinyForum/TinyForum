// job/clean_temp_files.go
package job

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/repository/attachment"
	driver "tiny-forum/internal/storage"
	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CleanTempFiles 周期性清理超过 24 小时的临时文件；ctx 取消时退出并释放 ticker。
func CleanTempFiles(ctx context.Context, db *gorm.DB, storage driver.StorageDriver, repo attachment.AttachmentRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Info("临时文件清理任务退出", zap.Error(ctx.Err()))
			return
		case <-ticker.C:
			cleanExpiredTempFiles(ctx, db, storage, repo)
		}
	}
}

// cleanExpiredTempFiles 执行一次过期临时文件清理
func cleanExpiredTempFiles(ctx context.Context, db *gorm.DB, storage driver.StorageDriver, repo attachment.AttachmentRepository) {
	var attachments []do.Attachment
	expireTime := time.Now().Add(-24 * time.Hour)
	if err := db.WithContext(ctx).
		Where("status = ? AND created_at < ?", do.StatusTemp, expireTime).
		Find(&attachments).Error; err != nil {
		logger.Error("查询过期临时文件失败", zap.Error(err))
		return
	}
	for _, att := range attachments {
		exists, statErr := storage.Exists(att.StoredPath)
		if statErr != nil {
			logger.Warn("检查过期临时文件状态失败",
				zap.String("path", att.StoredPath),
				zap.Error(statErr))
			continue
		}
		// 物理文件仍存在时先删除；已不存在则视为已清理，继续删库记录，
		// 避免残留孤儿记录被每小时反复重试。
		if exists {
			if err := storage.Delete(att.StoredPath); err != nil {
				logger.Warn("删除过期临时文件物理文件失败",
					zap.String("path", att.StoredPath),
					zap.Error(err))
				continue
			}
		}
		if err := repo.Delete(ctx, att.FileID); err != nil {
			logger.Warn("删除过期临时文件记录失败",
				zap.String("file_id", att.FileID),
				zap.Error(err))
		}
	}
}
