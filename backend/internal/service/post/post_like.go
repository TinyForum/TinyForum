package post

import (
	"tiny-forum/internal/model"
)

// Like 点赞帖子
func (s *postService) Like(userID, postID uint) error {
	if err := s.postRepo.AddLike(userID, postID); err != nil {
		return err
	}
	_ = s.postRepo.IncrLikeCount(postID, 1)
	_ = s.userRepo.AddScore(userID, 2)
	post, _ := s.postRepo.FindByID(postID)
	if post != nil && post.AuthorID != userID {
		s.notifSvc.Create(post.AuthorID, &userID, model.NotifyLike,
			"有人点赞了你的帖子《"+post.Title+"》", &postID, "post")
	}
	return nil
}

// Unlike 取消点赞帖子
func (s *postService) Unlike(userID, postID uint) error {
	if err := s.postRepo.RemoveLike(userID, postID); err != nil {
		return err
	}
	return s.postRepo.IncrLikeCount(postID, -1)
}
