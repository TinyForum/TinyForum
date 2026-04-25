package board

import (
	"errors"
)

func (s *boardService) DeletePost(boardID, postID, userID uint, isAdmin bool) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.BoardID != boardID {
		return errors.New("帖子不属于该板块")
	}
	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod && !isAdmin {
		return errors.New("无权限删除此帖子")
	}
	s.writeLog(userID, boardID, "delete_post", "post", postID, "版主删除")
	return s.postRepo.Delete(postID)
}

func (s *boardService) PinPost(boardID, postID uint, pin bool) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.BoardID != boardID {
		return errors.New("帖子不属于该板块")
	}
	return s.postRepo.TogglePinInBoard(postID, pin)
}
