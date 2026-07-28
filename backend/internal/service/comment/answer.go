package comment

import (
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
)

// MarkAsAnswer 标记/取消标记为答案
func (s *commentService) MarkAsAnswer(answerID, userID uint, isAdmin bool, isAnswer bool) error {
	answer, err := s.commentRepo.FindByAnswerID(answerID)
	if err != nil {
		return apperrors.ErrCommentNotFound
	}
	post, err := s.postRepo.FindQuestionByQuestionID(answer.QuestionID)
	if err != nil {
		return apperrors.ErrPostNotFound
	}
	if post.Creation.AuthorID != userID && !isAdmin {
		return apperrors.ErrInsufficientPermission
	}
	return s.commentRepo.MarkAsAnswer(answerID, isAnswer)
}

// UnacceptAnswer 取消接受答案（问题作者或管理员）
func (s *commentService) UnacceptAnswer(answerID, userID uint, isAdmin bool) error {
	answer, err := s.commentRepo.FindByAnswerID(answerID)
	if err != nil {
		return apperrors.ErrAnswerNotFound
	}
	post, err := s.postRepo.FindQuestionByQuestionID(answer.QuestionID)
	if err != nil {
		return apperrors.ErrQuestionNotFound
	}
	if post.Creation.Type != do.CreationTypeQuestion {
		return apperrors.ErrPostNotQuestion
	}
	if post.Creation.AuthorID != userID && !isAdmin {
		return apperrors.ErrInsufficientPermission // 权限不足，只有问题作者或管理员可以取消接受答案
	}
	if !answer.IsAccepted {
		return apperrors.ErrAnswerNotAccepted // 该答案未被接受
	}
	if err := s.commentRepo.UnacceptAnswer(answerID); err != nil {
		return err
	}
	// 可选：扣除积分并发送通知
	if post.Creation.AuthorID != userID {
		s.notifSvc.Create(answer.Reply.AuthorID, &userID, do.NotifyAcceptCancel,
			"你的答案在问题《"+post.Creation.Title+"》中被取消接受", &answer.QuestionID, "post")
	}
	return nil
}
