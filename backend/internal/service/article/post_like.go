package article

import (
	"tiny-forum/internal/model/do"
)

// Like 点赞帖子
func (s *articleService) Like(userID, postID uint) error {
	if err := s.postRepo.AddLike(userID, postID); err != nil {
		return err
	}

	err := s.postRepo.IncrLikeCount(postID, 1)
	if err != nil {
		return err
	}

	// 增加积分
	_ = s.userRepo.AddScore(userID, 2)
	post, _ := s.postRepo.FindByArticleID(postID)
	if post != nil && post.Creation.AuthorID != userID {
		s.notifSvc.Create(post.Creation.AuthorID, &userID, do.NotifyLike,
			"有人点赞了你的帖子《"+post.Creation.Title+"》", &postID, "post")
	}
	return nil
}

// Unlike 取消点赞帖子
func (s *articleService) Unlike(userID, postID uint) error {
	if err := s.postRepo.RemoveLike(userID, postID); err != nil {
		return err
	}
	return s.postRepo.IncrLikeCount(postID, -1)
}
