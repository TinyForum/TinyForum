package article

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
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

	// 异步记录点赞行为，供推荐系统使用
	go s.recSvc.RecordBehavior(userID, request.RecordBehaviorRequest{
		TargetID:     postID,
		TargetType:   "creation",
		BehaviorType: string(do.BehaviorLike),
		Value:        2.0,
	})

	return nil
}

// Unlike 取消点赞帖子
func (s *articleService) Unlike(userID, postID uint) error {
	if err := s.postRepo.RemoveLike(userID, postID); err != nil {
		return err
	}

	// 异步记录取消点赞行为
	go s.recSvc.RecordBehavior(userID, request.RecordBehaviorRequest{
		TargetID:     postID,
		TargetType:   "creation",
		BehaviorType: string(do.BehaviorUnlike),
		Value:        -1.0,
	})

	return s.postRepo.IncrLikeCount(postID, -1)
}
