package vote

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

// GetUserVote 获取用户对指定答案的投票类型
// 返回值：voteType（nil 表示未投票），error
// 适用场景：判断当前登录用户是否已投票，用于前端显示“已赞/已踩”状态。
func (r *voteRepository) GetUserVote(answerID, userID uint) (*do.AnswerVoteType, error) {
	logger.Debugf("[REPO] 获取用户投票: answerID=%d, userID=%d", answerID, userID)
	var vote do.AnswerVote
	err := r.db.Where("answer_id = ? AND user_id = ?", answerID, userID).
		Select("vote_type").
		First(&vote).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 无投票记录
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户投票失败: %w", err)
	}
	return vote.VoteType, nil // 注意返回指针
}

// GetAnswerVoteStats 获取答案的赞/踩统计（直接从 answers 表读取，性能优异）
// 要求：answers 表已增加 up_votes 和 down_votes 字段，并在投票事务中同步更新。
// 返回：upCount, downCount, error
func (r *voteRepository) GetAnswerVoteStats(answerID uint) (int, int, error) {
	var answer do.Answer
	err := r.db.Select("up_votes", "down_votes").Where("id = ?", answerID).First(&answer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, nil // 答案不存在，返回 0,0（或可返回错误）
		}
		return 0, 0, fmt.Errorf("查询答案投票统计失败: %w", err)
	}
	return answer.UpVotes, answer.DownVotes, nil
}

// GetVoteUsers 获取对某答案投了指定票型的用户列表（用于管理后台或审核）
// 注意：此方法会关联查询 answer_votes 和 users 表，数据量大时可能慢，请谨慎使用。
func (r *voteRepository) GetVoteUsers(answerID uint, voteType do.AnswerVoteType) ([]do.User, error) {
	var users []do.User
	err := r.db.Table("users").
		Joins("INNER JOIN answer_votes ON answer_votes.user_id = users.id").
		Where("answer_votes.answer_id = ? AND answer_votes.vote_type = ?", answerID, voteType).
		Find(&users).Error
	if err != nil {
		return nil, apperrors.ErrQueryVoteUserListFailed
	}
	return users, nil
}

// （可选）如果需要保留净票数方法，可以这样实现（但建议直接使用 GetAnswerVoteStats 计算差值）
// func (r *voteRepository) GetVoteCount(answerID uint) (int, error) {
//     up, down, err := r.GetAnswerVoteStats(answerID)
//     if err != nil {
//         return 0, err
//     }
//     return up - down, nil
// }
