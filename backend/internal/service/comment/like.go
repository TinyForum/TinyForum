package comment

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
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
	comment, _ := s.commentRepo.FindByCommentID(postID)
	if comment != nil && comment.Reply.AuthorID != userID {
		s.notifSvc.Create(comment.Reply.AuthorID, &userID, do.NotifyLike,
			"有人点赞了你的评论《"+comment.Reply.Content+"》", &postID, "post")
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
func (s *commentService) Unlike(userID, postID uint) error {
	if err := s.commentRepo.RemoveLike(userID, postID); err != nil {
		return err
	}

	// 异步记录取消点赞行为
	go s.recSvc.RecordBehavior(userID, request.RecordBehaviorRequest{
		TargetID:     postID,
		TargetType:   "creation",
		BehaviorType: string(do.BehaviorUnlike),
		Value:        -1.0,
	})

	return s.commentRepo.IncrLikeCount(postID, -1)
}
