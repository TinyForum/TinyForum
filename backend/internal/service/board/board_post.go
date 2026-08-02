package board

import (
	"context"
	apperrors "tiny-forum/pkg/errors"
)

func (s *boardService) DeletePost(boardID, postID, userID uint, isAdmin bool) error {
	post, err := s.postRepo.FindByArticleID(postID)
	if err != nil {
		return apperrors.ErrPostNotFound
	}
	if post.Creation.BoardID != boardID {
		return apperrors.ErrPostNotBelongToBoard
	}
	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod && !isAdmin {
		return apperrors.ErrInsufficientPermission
	}
	s.writeLog(userID, boardID, "delete_post", "post", postID, "版主删除")
	// 级联删除帖子关联的附件
	_ = s.attachmentSvc.DeleteByPostID(context.Background(), int64(postID))
	return s.postRepo.Delete(postID)
}

func (s *boardService) PinPost(boardID, postID uint, pin bool) error {
	post, err := s.postRepo.FindByArticleID(postID)
	if err != nil {
		return apperrors.ErrPostNotFound
	}
	if post.Creation.BoardID != boardID {
		return apperrors.ErrPostNotBelongToBoard
	}
	return s.postRepo.TogglePinInBoard(postID, pin)
}
