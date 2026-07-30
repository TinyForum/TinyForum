package comment

import (
	"fmt"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
)

// VoteAnswer 对回答进行投票（up/down），支持切换或取消（相同投票类型则取消）
func (s *commentService) VoteAnswer(answerID uint, userID uint, voteType do.AnswerVoteType) (*do.Answer, error) {
	// 1. 校验回答存在且不是自己的
	comment, err := s.commentRepo.FindByAnswerID(answerID)
	if err != nil {
		return nil, err
	}
	if comment.Reply.AuthorID == userID {
		return nil, apperrors.ErrCannotVoteSelfAnswer
	}

	// 2. 直接调用 Repository 的原子操作（内部有行锁，保证一致性）
	if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, voteType); err != nil {
		return nil, err
	}

	// 3. 重新获取最新评论（vote_count 已更新）
	updatedComment, err := s.commentRepo.FindByAnswerID(answerID)
	if err != nil {
		return nil, err
	}
	return updatedComment, nil
}

// RemoveVote 取消用户对回答的投票（无论当前是何类型）
func (s *commentService) RemoveVote(answerID uint, userID uint) (*do.Answer, error) {
	if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
		return nil, fmt.Errorf("取消投票失败: %w", err)
	}

	updatedComment, err := s.commentRepo.FindByAnswerID(answerID)
	if err != nil {
		return nil, fmt.Errorf("获取最新回答信息失败: %w", err)
	}
	return updatedComment, nil
}

// GetUserVoteStatus 获取用户对指定回答的投票状态
// 返回值：voteType（可能为 nil 表示未投票），error
func (s *commentService) GetUserVoteStatus(answerID uint, userID uint) (*do.AnswerVoteType, error) {
	logger.Debugf("[SERVICE] 获取用户投票状态: answerID=%d, userID=%d")
	voteType, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户投票状态失败: %w", err)
	}
	return voteType, nil
}
