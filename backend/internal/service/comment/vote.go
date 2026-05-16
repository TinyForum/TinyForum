package comment

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"
)

// VoteAnswer 对回答进行投票（up/down），支持切换或取消（相同投票类型则取消）
func (s *commentService) VoteAnswer(answerID uint, userID uint, voteType do.AnswerVoteType) (*do.Comment, error) {
	// 1. 校验回答存在且为答案类型
	comment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, fmt.Errorf("回答不存在: %w", err)
	}
	if !comment.IsAnswer {
		return nil, errors.New("该评论不是回答，无法投票")
	}
	if comment.AuthorID == userID {
		return nil, errors.New("不能给自己的回答投票")
	}

	// 2. 获取用户当前投票状态
	currentVote, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户投票状态失败: %w", err)
	}

	// 3. 根据当前状态决定操作
	switch {
	case currentVote != nil && *currentVote == voteType:
		// 相同投票 → 取消投票
		if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
			return nil, fmt.Errorf("取消投票失败: %w", err)
		}
	case currentVote == nil:
		// 未投票 → 新增投票
		if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, voteType); err != nil {
			return nil, fmt.Errorf("创建投票失败: %w", err)
		}
	default:
		// 已投票但类型不同 → 更新投票（先删后加或原子更新）
		if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, voteType); err != nil {
			return nil, fmt.Errorf("更新投票失败: %w", err)
		}
	}

	// 4. 重新获取最新评论（包含更新后的 vote_count）
	updatedComment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, fmt.Errorf("获取最新回答信息失败: %w", err)
	}
	return updatedComment, nil
}

// RemoveVote 取消用户对回答的投票（无论当前是何类型）
func (s *commentService) RemoveVote(answerID uint, userID uint) (*do.Comment, error) {
	// 1. 校验回答存在且为答案类型
	comment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, fmt.Errorf("回答不存在: %w", err)
	}
	if !comment.IsAnswer {
		return nil, errors.New("该评论不是回答，无法取消投票")
	}

	// 2. 检查是否有投票记录
	currentVote, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户投票状态失败: %w", err)
	}
	if currentVote == nil {
		return nil, errors.New("尚未投票，无法取消")
	}

	// 3. 执行取消
	if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
		return nil, fmt.Errorf("取消投票失败: %w", err)
	}

	// 4. 返回最新评论
	updatedComment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, fmt.Errorf("获取最新回答信息失败: %w", err)
	}
	return updatedComment, nil
}

// GetUserVoteStatus 获取用户对指定回答的投票状态
// 返回值：voteType（可能为 nil 表示未投票），error
func (s *commentService) GetUserVoteStatus(answerID uint, userID uint) (*do.AnswerVoteType, error) {
	voteType, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户投票状态失败: %w", err)
	}
	return voteType, nil
}
