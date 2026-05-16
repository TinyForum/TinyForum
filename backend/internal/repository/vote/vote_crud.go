package vote

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// voteTypeToDelta 将投票类型映射为 vote_count 的变化量
func voteTypeToDelta(vt do.AnswerVoteType) int {
	if vt == do.AnswerVoteTypeUp {
		return 1
	}
	return -1 // down
}

// CreateOrUpdateVote 创建或更新投票（支持 up/down）
func (r *voteRepository) CreateOrUpdateVote(commentID, userID uint, voteType do.AnswerVoteType) error {
	// // 1. 解析并校验投票类型
	// voteType, err := do.ParseAnswerVoteType(voteTypeStr)
	// if err != nil {
	// 	return fmt.Errorf("无效的投票类型: %w", err)
	// }

	// 2. 使用 GORM 事务（自动处理 commit/rollback）
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingVote do.AnswerVote

		// 2.1 查询是否已存在投票
		err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
			First(&existingVote).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询已有投票失败: %w", err)
		}

		newDelta := voteTypeToDelta(voteType)

		if existingVote.ID > 0 {
			// 更新已存在的投票
			oldDelta := voteTypeToDelta(*existingVote.VoteType) // 注意指针解引用
			diff := newDelta - oldDelta

			// 更新投票记录
			if err := tx.Model(&existingVote).Update("vote_type", voteType).Error; err != nil {
				return fmt.Errorf("更新投票失败: %w", err)
			}

			// 更新评论的 vote_count（如果有变化）
			if diff != 0 {
				if err := tx.Model(&do.Comment{}).
					Where("id = ?", commentID).
					UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error; err != nil {
					return fmt.Errorf("更新评论投票计数失败: %w", err)
				}
			}
		} else {
			// 新增投票
			vote := &do.AnswerVote{
				UserID:    userID,
				CommentID: commentID,
				VoteType:  &voteType, // 若字段为指针，需取地址
			}
			if err := tx.Create(vote).Error; err != nil {
				return fmt.Errorf("创建投票失败: %w", err)
			}

			// 增加评论的 vote_count
			if err := tx.Model(&do.Comment{}).
				Where("id = ?", commentID).
				UpdateColumn("vote_count", gorm.Expr("vote_count + ?", newDelta)).Error; err != nil {
				return fmt.Errorf("更新评论投票计数失败: %w", err)
			}
		}

		return nil
	})
}

// RemoveVote 删除投票（取消 up/down）
func (r *voteRepository) RemoveVote(commentID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var vote do.AnswerVote

		// 1. 查询投票
		if err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
			First(&vote).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil // 没有投票，视为已删除成功
			}
			return fmt.Errorf("查询投票失败: %w", err)
		}

		// 2. 计算对 vote_count 的影响量
		delta := voteTypeToDelta(*vote.VoteType)

		// 3. 删除投票记录
		if err := tx.Delete(&vote).Error; err != nil {
			return fmt.Errorf("删除投票失败: %w", err)
		}

		// 4. 减少评论的 vote_count
		if err := tx.Model(&do.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("vote_count", gorm.Expr("vote_count - ?", delta)).Error; err != nil {
			return fmt.Errorf("更新评论投票计数失败: %w", err)
		}

		return nil
	})
}
