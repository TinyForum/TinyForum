package comment

import (
	"errors"
	apperrors "tiny-forum/pkg/errors"
)

// Delete 删除普通评论
func (s *commentService) Delete(commentID, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByCommentID(commentID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	if comment.Reply.AuthorID != userID && !isAdmin {
		return errors.New("无权限删除此评论")
	}
	return s.commentRepo.Delete(commentID)
}

// DeleteAnswer 删除回答（权限：管理员、作者、问题作者）
func (s *commentService) DeleteAnswer(commentID, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByAnswerID(commentID)
	if err != nil {
		return errors.New("回答不存在")
	}

	if isAdmin {
		return s.commentRepo.Delete(commentID)
	}
	if comment.Reply.AuthorID == userID {
		return s.commentRepo.Delete(commentID)
	}
	post, err := s.postRepo.FindQuestionByQuestionID(comment.QuestionID)
	if err == nil && post.Creation.AuthorID == userID {
		return s.commentRepo.Delete(commentID)
	}
	return errors.New("无权限删除此回答")
}
