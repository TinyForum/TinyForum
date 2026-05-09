// job/clean_temp_files.go
package job

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/repository/attachment"
	driver "tiny-forum/internal/storage"

	"gorm.io/gorm"
)

// 清理临时文件的定时任务
func CleanTempFiles(db *gorm.DB, storage driver.StorageDriver, repo attachment.AttachmentRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			ctx := context.Background()
			var attachments []do.Attachment
			expireTime := time.Now().Add(-24 * time.Hour)
			db.Where("status = ? AND created_at < ?", do.StatusTemp, expireTime).Find(&attachments)
			for _, att := range attachments {
				_ = storage.Delete(att.StoredPath)
				_ = repo.Delete(ctx, att.FileID)
			}
		}
	}()
}
