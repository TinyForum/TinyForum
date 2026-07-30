package vote

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// updateAnswerVoteCount 原子更新回答的 up_votes 或 down_votes
// voteType: 投票类型（up/down），delta: +1 或 -1
func updateAnswerVoteCount(tx *gorm.DB, answerID uint, voteType do.AnswerVoteType, delta int) error {
	if delta == 0 {
		logger.Debugf("updateAnswerVoteCount: delta=0, skip for answerID=%d", answerID)
		return nil
	}

	// 确定要更新的字段名
	field := "up_votes"
	if voteType == do.AnswerVoteTypeDown {
		field = "down_votes"
	}

	logger.Debugf("updateAnswerVoteCount: answerID=%d, field=%s, delta=%d", answerID, field, delta)

	// 原子更新
	if err := tx.Model(&do.Answer{}).
		Where("id = ?", answerID).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta)).Error; err != nil {
		return fmt.Errorf("更新回答投票计数失败: %w", err)
	}
	return nil
}

// CreateOrUpdateVote 创建或更新投票（支持 up/down）
func (r *voteRepository) CreateOrUpdateVote(answerID, userID uint, voteType do.AnswerVoteType) error {
	logger.Infof("CreateOrUpdateVote: answerID=%d, userID=%d, voteType=%s", answerID, userID, voteType)

	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing do.AnswerVote
		err := tx.Unscoped().
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("answer_id = ? AND user_id = ?", answerID, userID).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrQueryVoteRecordFailed
		}

		// 场景1：记录不存在 → 创建新投票
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debugf("CreateOrUpdateVote: no existing record, creating new vote")
			vote := &do.AnswerVote{
				UserID:   userID,
				AnswerID: answerID,
				VoteType: &voteType,
			}
			if err := tx.Create(vote).Error; err != nil {
				return apperrors.ErrCreateVoteFailed
			}
			// 增加对应类型的计数
			return updateAnswerVoteCount(tx, answerID, voteType, 1)
		}

		// 记录存在
		oldVoteType := *existing.VoteType
		isDeleted := existing.DeletedAt.Valid && !existing.DeletedAt.Time.IsZero()
		logger.Debugf("CreateOrUpdateVote: existing vote id=%d, oldType=%s, newType=%s, isDeleted=%v",
			existing.ID, oldVoteType, voteType, isDeleted)

		// 准备更新字段
		updateMap := map[string]interface{}{"vote_type": voteType}
		if isDeleted {
			updateMap["deleted_at"] = nil // 恢复软删除
		}

		if err := tx.Unscoped().Model(&existing).Updates(updateMap).Error; err != nil {
			return apperrors.ErrUpdateVoteRecordFailed
		}

		// 计算需要调整的计数
		if isDeleted {
			// 恢复记录，只需增加新投票类型的计数
			return updateAnswerVoteCount(tx, answerID, voteType, 1)
		} else {
			// 更改类型：先减掉旧类型的计数，再增加新类型的计数
			if oldVoteType == voteType {
				// 类型未变，无需更新计数（但可能恢复软删，已经处理）
				return nil
			}
			// 减去旧类型
			if err := updateAnswerVoteCount(tx, answerID, oldVoteType, -1); err != nil {
				return err
			}
			// 增加新类型
			return updateAnswerVoteCount(tx, answerID, voteType, 1)
		}
	})
}

// RemoveVote 取消投票（软删除）
func (r *voteRepository) RemoveVote(answerID, userID uint) error {
	logger.Infof("RemoveVote: answerID=%d, userID=%d", answerID, userID)
	return r.db.Transaction(func(tx *gorm.DB) error {
		var vote do.AnswerVote
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("answer_id = ? AND user_id = ?", answerID, userID).
			First(&vote).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Debugf("RemoveVote: no active vote found, nothing to do")
				return nil
			}
			return apperrors.ErrQueryVoteRecordFailed
		}

		voteType := *vote.VoteType
		logger.Debugf("RemoveVote: found vote id=%d, type=%s", vote.ID, voteType)

		// 软删除
		if err := tx.Delete(&vote).Error; err != nil {
			return apperrors.ErrDeleteVoteFailed
		}

		// 减去对应类型的计数
		return updateAnswerVoteCount(tx, answerID, voteType, -1)
	})
}
