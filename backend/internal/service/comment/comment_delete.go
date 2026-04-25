package comment

import (
	"errors"
)

// Delete 删除普通评论
func (s *commentService) Delete(commentID, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	if comment.AuthorID != userID && !isAdmin {
		return errors.New("无权限删除此评论")
	}
	return s.commentRepo.Delete(commentID)
}

// DeleteAnswer 删除回答（权限：管理员、作者、问题作者）
func (s *commentService) DeleteAnswer(commentID, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return errors.New("该评论不是回答")
	}
	if isAdmin {
		return s.commentRepo.Delete(commentID)
	}
	if comment.AuthorID == userID {
		return s.commentRepo.Delete(commentID)
	}
	post, err := s.postRepo.FindByID(comment.PostID)
	if err == nil && post.AuthorID == userID {
		return s.commentRepo.Delete(commentID)
	}
	return errors.New("无权限删除此回答")
}
