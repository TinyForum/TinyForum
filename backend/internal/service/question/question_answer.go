package question

import (
	"errors"
	"fmt"

	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
)

type VoteAnswerInput struct {
	CommentID uint   `json:"comment_id" binding:"required"`
	VoteType  string `json:"vote_type" binding:"required,oneof=up down"`
}

type VoteAnswerResult struct {
	VoteType  string `json:"vote_type"`
	VoteCount int    `json:"vote_count"`
	Action    string `json:"action"`
}

// AcceptAnswer 采纳答案
func (s *questionService) AcceptAnswer(postID, commentID uint, userID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return fmt.Errorf("%w: 帖子不存在", apperrors.ErrPostNotFound)
	}
	if post.AuthorID != userID {
		return fmt.Errorf("%w: 只有发帖人可以采纳答案", apperrors.ErrAcceptForbidden)
	}
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return fmt.Errorf("%w: 回答不存在", apperrors.ErrAnswerNotFound)
	}
	if !comment.IsAnswer {
		return errors.New("该评论不是回答")
	}
	question, err := s.questionRepo.FindByPostID(postID)
	if err != nil {
		return fmt.Errorf("%w: 问答信息不存在", apperrors.ErrQuestionNotFound)
	}
	if question.AcceptedAnswerID != nil {
		return errors.New("已经采纳过答案了")
	}
	if err := s.questionRepo.SetAcceptedAnswer(postID, commentID); err != nil {
		return err
	}
	if err := s.commentRepo.MarkAsAccepted(commentID); err != nil {
		return err
	}
	if question.RewardScore > 0 {
		s.userRepo.AddScore(comment.AuthorID, question.RewardScore)
	}
	s.notifSvc.Create(comment.AuthorID, &userID, do.NotifySystem,
		"你的回答被采纳为最佳答案", &postID, "post")
	return nil
}

// VoteAnswer 投票回答
func (s *questionService) VoteAnswer(userID uint, input VoteAnswerInput) (*VoteAnswerResult, error) {
	comment, err := s.commentRepo.FindByID(input.CommentID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return nil, errors.New("只能对回答进行投票")
	}
	if comment.AuthorID == userID {
		return nil, errors.New("不能给自己的答案投票")
	}
	existingVote, _ := s.questionRepo.FindAnswerVote(userID, input.CommentID)
	var result VoteAnswerResult
	var action string
	if existingVote != nil && existingVote.ID != 0 {
		if existingVote.VoteType == input.VoteType {
			if err := s.questionRepo.DeleteAnswerVote(userID, input.CommentID); err != nil {
				return nil, err
			}
			action = "removed"
			result.VoteType = ""
		} else {
			existingVote.VoteType = input.VoteType
			if err := s.questionRepo.UpdateAnswerVote(existingVote); err != nil {
				return nil, err
			}
			action = "updated"
			result.VoteType = input.VoteType
		}
	} else {
		vote := &do.AnswerVote{
			UserID:    userID,
			CommentID: input.CommentID,
			VoteType:  input.VoteType,
		}
		if err := s.questionRepo.CreateAnswerVote(vote); err != nil {
			return nil, err
		}
		action = "added"
		result.VoteType = input.VoteType
	}
	voteCount, _ := s.questionRepo.GetAnswerVoteCount(input.CommentID)
	s.commentRepo.UpdateVoteCount(input.CommentID, voteCount)
	result.VoteCount = voteCount
	result.Action = action
	if action != "removed" {
		s.notifSvc.Create(comment.AuthorID, &userID, do.NotifyLike,
			"有人给你的答案投票了", &input.CommentID, "comment")
	}
	return &result, nil
}

// GetAnswerVoteStatus 获取用户对答案的投票状态
func (s *questionService) GetAnswerVoteStatus(userID, commentID uint) (map[string]interface{}, error) {
	vote, err := s.questionRepo.FindAnswerVote(userID, commentID)
	if err != nil {
		return map[string]interface{}{
			"has_voted":  false,
			"vote_type":  "",
			"vote_count": 0,
		}, nil
	}
	voteCount, _ := s.questionRepo.GetAnswerVoteCount(commentID)
	return map[string]interface{}{
		"has_voted":  true,
		"vote_type":  vote.VoteType,
		"vote_count": voteCount,
	}, nil
}

// GetQuestionWithAnswers 获取问题及其回答（分页）
func (s *questionService) GetQuestionWithAnswers(postID uint, page, pageSize int) (*do.Question, []do.Comment, int64, error) {
	question, err := s.questionRepo.FindByPostID(postID)
	if err != nil {
		return nil, nil, 0, err
	}
	answers, total, err := s.commentRepo.GetAnswersByPostID(postID, pageSize, (page-1)*pageSize)
	return question, answers, total, err
}
