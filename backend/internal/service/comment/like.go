package comment

import (
	"tiny-forum/internal/model/do"
)

// Like 点赞帖子
func (s *commentService) Like(userID, postID uint) error {
	if err := s.commentRepo.AddLike(userID, postID); err != nil {
		return err
	}

	err := s.postRepo.IncrLikeCount(postID, 1)
	if err != nil {
		return err
	}

	// 增加积分
	_ = s.userRepo.AddScore(userID, 2)
	comment, _ := s.commentRepo.FindByID(postID)
	if comment != nil && comment.AuthorID != userID {
		s.notifSvc.Create(comment.AuthorID, &userID, do.NotifyLike,
			"有人点赞了你的评论《"+comment.Content+"》", &postID, "post")
	}
	return nil
}

// Unlike 取消点赞帖子
func (s *commentService) Unlike(userID, postID uint) error {
	if err := s.commentRepo.RemoveLike(userID, postID); err != nil {
		return err
	}
	return s.commentRepo.IncrLikeCount(postID, -1)
}
