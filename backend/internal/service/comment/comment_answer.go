package comment

import (
	"errors"
	"tiny-forum/internal/model/po"
)

// MarkAsAnswer 标记/取消标记为答案
func (s *commentService) MarkAsAnswer(commentID, userID uint, isAdmin bool, isAnswer bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	post, err := s.postRepo.FindByID(comment.PostID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.AuthorID != userID && !isAdmin {
		return errors.New("无权限操作")
	}
	return s.commentRepo.MarkAsAnswer(commentID, isAnswer)
}

// UnacceptAnswer 取消接受答案（问题作者或管理员）
func (s *commentService) UnacceptAnswer(answerID, userID uint, isAdmin bool) error {
	answer, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return errors.New("回答不存在")
	}
	if !answer.IsAnswer {
		return errors.New("该评论不是回答")
	}
	post, err := s.postRepo.FindByID(answer.PostID)
	if err != nil {
		return errors.New("问题不存在")
	}
	if post.Type != "question" {
		return errors.New("该帖子不是问答类型")
	}
	if post.AuthorID != userID && !isAdmin {
		return errors.New("没有权限操作，只有问题作者可以取消接受答案")
	}
	if !answer.IsAccepted {
		return errors.New("该回答未被接受为答案")
	}
	if err := s.commentRepo.UnacceptAnswer(answerID); err != nil {
		return err
	}
	// 可选：扣除积分并发送通知
	if post.AuthorID != userID {
		s.notifSvc.Create(answer.AuthorID, &userID, po.NotifyAcceptCancel,
			"你的答案在问题《"+post.Title+"》中被取消接受", &answer.PostID, "post")
	}
	return nil
}
