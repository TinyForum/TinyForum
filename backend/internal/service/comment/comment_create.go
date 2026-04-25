package comment

import (
	"errors"

	"tiny-forum/internal/model"
)

type CreateCommentInput struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=2000"`
	ParentID *uint  `json:"parent_id"`
}

// Create 创建普通评论
func (s *commentService) Create(authorID uint, input CreateCommentInput) (*model.Comment, error) {
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}

	if input.ParentID != nil && *input.ParentID != 0 {
		if err := s.commentRepo.ValidateParentComment(*input.ParentID, input.PostID); err != nil {
			return nil, err
		}
	}
	comment := &model.Comment{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: authorID,
		ParentID: input.ParentID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	_ = s.userRepo.AddScore(authorID, 3)

	if post.AuthorID != authorID {
		s.notifSvc.Create(post.AuthorID, &authorID, model.NotifyComment,
			"有人评论了你的帖子《"+post.Title+"》", &input.PostID, "post")
	}

	if input.ParentID != nil {
		parent, err := s.commentRepo.FindByID(*input.ParentID)
		if err == nil && parent.AuthorID != authorID {
			s.notifSvc.Create(parent.AuthorID, &authorID, model.NotifyReply,
				"有人回复了你的评论", input.ParentID, "comment")
		}
	}

	return s.commentRepo.FindByID(comment.ID)
}

// CreateAnswer 创建回答（仅限问答帖）
func (s *commentService) CreateAnswer(authorID uint, input CreateCommentInput) (*model.Comment, error) {
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}
	if post.Type != "question" {
		return nil, errors.New("该帖子不是问答类型，请使用普通评论")
	}

	comment := &model.Comment{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: authorID,
		ParentID: input.ParentID,
		IsAnswer: true,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	_ = s.userRepo.AddScore(authorID, 2)

	if post.AuthorID != authorID {
		s.notifSvc.Create(post.AuthorID, &authorID, model.NotifyComment,
			"有人回答了你的问题《"+post.Title+"》", &input.PostID, "post")
	}

	return s.commentRepo.FindByID(comment.ID)
}
