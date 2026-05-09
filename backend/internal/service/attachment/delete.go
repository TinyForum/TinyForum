// package attachment

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"tiny-forum/pkg/logger"
// )

// // 删除文件
// func (s *service) DeleteFile(ctx context.Context, userID int64, fileID string) error {
// 	att, err := s.repo.GetByFileID(ctx, fileID)
// 	if err != nil {
// 		return err
// 	}

// 	// 检查权限：只有文件所有者或管理员可以删除
// 	if att.UserID != userID {
// 		// TODO: 添加管理员权限检查
// 		return fmt.Errorf("无权删除此文件")
// 	}

// 	// 删除数据库记录
// 	if err := s.repo.Delete(ctx, fileID); err != nil {
// 		return err
// 	}

// 	// 删除物理文件（可选：保留一段时间或异步删除）
// 	if err := os.Remove(att.StoredPath); err != nil {
// 		logger.Warnf("删除物理文件失败: %v", err)
// 	}

//		return nil
//	}
package attachment