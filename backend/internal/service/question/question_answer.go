package question

import (
	"context"
	"errors"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// AcceptAnswer 采纳答案
func (s *questionService) AcceptAnswer(questionID, answerID uint, userID uint) error {
	question, err := s.postRepo.FindQuestionByQuestionID(questionID)
	if err != nil {
		return apperrors.ErrPostNotFound
	}
	if question.Creation.AuthorID != userID {
		return apperrors.ErrAcceptForbidden
	}
	answer, err := s.commentRepo.FindByAnswerID(answerID)
	if err != nil {
		return apperrors.ErrAnswerNotFound
	}

	if question.AcceptedAnswerID != nil {
		if *question.AcceptedAnswerID != answerID {
			return apperrors.ErrAlreadyAcceptedAnswer
		}
		if answer.IsAccepted {
			return apperrors.ErrAlreadyAcceptedAnswer
		}
	}

	if err := s.txManager.ExecuteInTransaction(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Model(&do.Question{}).Where("id = ?", questionID).
			Updates(map[string]interface{}{
				"accepted_answer_id": answerID,
			}).Error; err != nil {
			return err
		}
		return tx.Model(&do.Answer{}).Where("id = ?", answerID).
			Update("is_accepted", true).Error
	}); err != nil {
		return err
	}
	if question.RewardScore > 0 {
		s.userRepo.AddScore(answer.Reply.AuthorID, question.RewardScore)
	}
	s.notifSvc.Create(answer.Reply.AuthorID, &userID, do.NotifySystem,
		"你的回答被采纳为最佳答案", &questionID, "post")
	return nil
}

// VoteAnswer 投票回答
func (s *questionService) VoteAnswer(userID uint, input request.VoteAnswerRequest) (*vo.VoteAnswerVO, error) {
	comment, err := s.commentRepo.FindByAnswerID(input.CommentID)
	if err != nil {
		return nil, apperrors.ErrAnswerNotFound
	}
	// if !comment.IsAnswer {
	// 	return nil, errors.New("只能对回答进行投票")
	// }
	if comment.Reply.AuthorID == userID {
		return nil, errors.New("不能给自己的答案投票")
	}
	existingVote, _ := s.questionRepo.FindAnswerVote(userID, input.CommentID)
	var result vo.VoteAnswerVO
	var action string
	// 如果用户已经投过票，则更新投票类型
	if existingVote != nil && existingVote.ID != 0 {
		if existingVote.VoteType == input.VoteType {
			if err := s.questionRepo.DeleteAnswerVote(userID, input.CommentID); err != nil {
				return nil, err
			}
			// 如果投票类型相同，则删除投票记录
			action = "removed"
			result.VoteType = nil
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
			UserID:   userID,
			AnswerID: input.CommentID,
			VoteType: input.VoteType,
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
		s.notifSvc.Create(comment.Reply.AuthorID, &userID, do.NotifyLike,
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
func (s *questionService) GetAnswersList(questionID uint, page, pageSize int) ([]do.Answer, int64, error) {
	// question, err := s.questionRepo.FindByQuestionID(questionID)
	// if err != nil {
	// 	return nil, 0, err
	// }
	answers, total, err := s.commentRepo.GetAnswersByQuestionID(questionID, pageSize, (page-1)*pageSize)
	return answers, total, err
}
