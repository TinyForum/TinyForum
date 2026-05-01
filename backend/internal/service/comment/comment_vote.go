package comment

import (
	"errors"

	"tiny-forum/internal/model/po"
)

// VoteAnswer 投票回答（支持 up/down）
func (s *commentService) VoteAnswer(answerID uint, userID uint, voteType int) (*po.Comment, error) {
	comment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return nil, errors.New("该评论不是回答")
	}
	if comment.AuthorID == userID {
		return nil, errors.New("不能给自己的回答投票")
	}

	currentVote, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, err
	}

	switch currentVote {
	case voteType:
		// 相同投票：取消投票
		if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
			return nil, err
		}
	case 0:
		// 未投票：创建新投票
		if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, voteType); err != nil {
			return nil, err
		}
	default:
		// 改变投票：更新现有投票
		if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, voteType); err != nil {
			return nil, err
		}
	}

	return s.commentRepo.FindByID(answerID)
}

// RemoveVote 取消投票
func (s *commentService) RemoveVote(answerID uint, userID uint) (*po.Comment, error) {
	comment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return nil, errors.New("该评论不是回答")
	}
	currentVote, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, err
	}
	if currentVote == 0 {
		return nil, errors.New("尚未投票，无法取消")
	}
	if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
		return nil, err
	}
	return s.commentRepo.FindByID(answerID)
}

// GetUserVoteStatus 获取用户对指定答案的投票状态
func (s *commentService) GetUserVoteStatus(answerID uint, userID uint) (int, error) {
	return s.voteRepo.GetUserVote(answerID, userID)
}
