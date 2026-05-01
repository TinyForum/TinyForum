package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

func (s *authService) DeleteAccount(ctx context.Context, userID uint, input DeleteAccountInput) error {
	if input.Confirm != "DELETE" {
		return errors.New("请确认删除操作")
	}

	// 通过 repo 软删除
	if err := s.authRepo.SoftDelete(ctx, userID); err != nil {
		return errors.New("删除账户失败")
	}

	return nil
}

// internal/service/auth/account.go

// ScheduleDeletion 标记账户为待删除（软删除）
func (s *authService) ScheduleDeletion(ctx context.Context, userID uint, input DeleteAccountInput) error {
	if input.Confirm != "DELETE" {
		return errors.New("请确认删除操作")
	}

	// 可选：验证密码
	if input.Password != "" {
		// if err := s.authRepo.VerifyPassword(ctx, userID, input.Password); err != nil {
		//     return errors.New("密码错误")
		// }
	}

	// 软删除用户（设置 deleted_at）
	return s.authRepo.SoftDelete(ctx, userID)
}

// CancelDeletion 取消注销，恢复账户
func (s *authService) CancelDeletion(ctx context.Context, userID uint) error {
	// 恢复账户：将 deleted_at 设置为 NULL
	return s.authRepo.Restore(ctx, userID)
}

// ConfirmDeletion 永久删除账户（硬删除）
func (s *authService) ConfirmDeletion(ctx context.Context, userID uint) error {
	return s.txManager.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
		// ========== 第一组：直接关联 user_id 的表 ==========

		// 1. 通知
		if err := tx.Exec("DELETE FROM notifications WHERE user_id = ? OR sender_id = ?", userID, userID).Error; err != nil {
			return fmt.Errorf("删除通知失败: %w", err)
		}

		// if err := tx.Where("user_id = ? OR sender_id = ?", userID, userID).Delete(&po.Notification{}).Error; err != nil {
		// 	return fmt.Errorf("删除通知失败: %w", err)
		// }
		// 2. 点赞
		if err := tx.Where("user_id = ?", userID).Delete(&po.Like{}).Error; err != nil {
			return fmt.Errorf("删除点赞失败: %w", err)
		}

		// 3. 投票
		if err := tx.Where("user_id = ?", userID).Delete(&po.Vote{}).Error; err != nil {
			return fmt.Errorf("删除投票失败: %w", err)
		}

		// 4. 回答投票
		if err := tx.Where("user_id = ?", userID).Delete(&po.AnswerVote{}).Error; err != nil {
			return fmt.Errorf("删除回答投票失败: %w", err)
		}

		// 5. 关注关系（作为关注者）
		if err := tx.Where("follower_id = ?", userID).Delete(&po.Follow{}).Error; err != nil {
			return fmt.Errorf("删除关注关系失败: %w", err)
		}

		// 6. 签到记录
		if err := tx.Where("user_id = ?", userID).Delete(&po.SignIn{}).Error; err != nil {
			return fmt.Errorf("删除签到记录失败: %w", err)
		}

		// 7. 版主申请
		if err := tx.Where("user_id = ? OR reviewer_id = ?", userID, userID).Delete(&po.ModeratorApplication{}).Error; err != nil {
			return fmt.Errorf("删除版主申请失败: %w", err)
		}

		// 8. 版主记录
		if err := tx.Where("moderator_id = ?", userID).Delete(&po.ModeratorLog{}).Error; err != nil {
			return fmt.Errorf("删除版主日志失败: %w", err)
		}

		// 9. 版主权限
		if err := tx.Where("user_id = ?", userID).Delete(&po.Moderator{}).Error; err != nil {
			return fmt.Errorf("删除版主权限失败: %w", err)
		}

		// 10. 举报（作为举报人或处理人）
		if err := tx.Where("reporter_id = ? OR handler_id = ?", userID, userID).Delete(&po.Report{}).Error; err != nil {
			return fmt.Errorf("删除举报记录失败: %w", err)
		}

		// 11. 时间线事件（作为用户或触发者）
		if err := tx.Where("user_id = ? OR actor_id = ?", userID, userID).Delete(&po.TimelineEvent{}).Error; err != nil {
			return fmt.Errorf("删除时间线事件失败: %w", err)
		}

		// 12. 时间线订阅
		if err := tx.Where("subscriber_id = ? OR target_user_id = ?", userID, userID).Delete(&po.TimelineSubscription{}).Error; err != nil {
			return fmt.Errorf("删除时间线订阅失败: %w", err)
		}

		// 13. 用户时间线
		if err := tx.Where("user_id = ?", userID).Delete(&po.UserTimeline{}).Error; err != nil {
			return fmt.Errorf("删除用户时间线失败: %w", err)
		}

		// 14. 主题关注
		if err := tx.Where("user_id = ?", userID).Delete(&po.TopicFollow{}).Error; err != nil {
			return fmt.Errorf("删除主题关注失败: %w", err)
		}

		// 15. 主题帖子关联
		if err := tx.Where("added_by = ?", userID).Delete(&po.TopicPost{}).Error; err != nil {
			return fmt.Errorf("删除主题帖子关联失败: %w", err)
		}

		// 16. 主题（创建者）
		if err := tx.Where("creator_id = ?", userID).Delete(&po.Topic{}).Error; err != nil {
			return fmt.Errorf("删除主题失败: %w", err)
		}

		// 17. 版块封禁
		if err := tx.Where("user_id = ? OR banned_by = ?", userID, userID).Delete(&po.BoardBan{}).Error; err != nil {
			return fmt.Errorf("删除版块封禁失败: %w", err)
		}

		// 18. 公告（创建者或更新者）
		if err := tx.Where("created_by = ? OR updated_by = ?", userID, userID).Delete(&po.Announcement{}).Error; err != nil {
			return fmt.Errorf("删除公告失败: %w", err)
		}

		// 19. 邀请关系
		if err := tx.Model(&po.User{}).Where("invited_by_id = ?", userID).Update("invited_by_id", nil).Error; err != nil {
			return fmt.Errorf("更新邀请关系失败: %w", err)
		}

		// ========== 第二组：处理需要保留但置空的表 ==========

		// 20. 评论（置空 author_id）
		if err := tx.Model(&po.Comment{}).Where("author_id = ?", userID).Update("author_id", nil).Error; err != nil {
			return fmt.Errorf("更新评论失败: %w", err)
		}

		// 21. 帖子（置空 author_id）
		if err := tx.Model(&po.Post{}).Where("author_id = ?", userID).Update("author_id", nil).Error; err != nil {
			return fmt.Errorf("更新帖子失败: %w", err)
		}

		// 22. 问题（置空 accepted_answer_id，如果答案属于该用户）
		// 先查出该用户的答案，然后置空
		if err := tx.Model(&po.Question{}).
			Where("accepted_answer_id IN (?)",
				tx.Table("comments").Select("id").Where("author_id = ?", userID)).
			Update("accepted_answer_id", nil).Error; err != nil {
			return fmt.Errorf("更新问题失败: %w", err)
		}

		// ========== 第三组：最后删除用户 ==========

		// 23. 删除用户
		if err := tx.Unscoped().Delete(&po.User{}, userID).Error; err != nil {
			return fmt.Errorf("删除用户失败: %w", err)
		}

		return nil
	})
}

// GetDeletionStatus 获取账户删除状态
func (s *authService) GetDeletionStatus(ctx context.Context, userID uint) (*DeletionStatus, error) {
	// 获取用户信息（包括已软删除的）
	user, err := s.authRepo.GetUserWithDeleted(ctx, userID)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}

	status := &DeletionStatus{
		IsDeleted:  false,
		CanRestore: false,
	}

	// 检查是否已软删除
	if user.DeletedAt.Valid {
		status.IsDeleted = true
		status.DeletedAt = &user.DeletedAt.Time

		// 计算剩余可恢复天数（例如30天内可恢复）
		daysSinceDeleted := int(time.Since(user.DeletedAt.Time).Hours() / 24)
		remainingDays := 30 - daysSinceDeleted

		if remainingDays > 0 {
			status.CanRestore = true
			status.RemainingDays = remainingDays
		}
	}

	return status, nil
}
