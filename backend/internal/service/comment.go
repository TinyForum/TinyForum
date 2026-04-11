package service

import (
	"errors"

	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
	userRepo    *repository.UserRepository
	notifSvc    *NotificationService
}

func NewCommentService(
	commentRepo *repository.CommentRepository,
	postRepo *repository.PostRepository,
	userRepo *repository.UserRepository,
	notifSvc *NotificationService,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifSvc:    notifSvc,
	}
}

type CreateCommentInput struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=2000"`
	ParentID *uint  `json:"parent_id"`
}

func (s *CommentService) Create(authorID uint, input CreateCommentInput) (*model.Comment, error) {
	// Validate post exists
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return nil, errors.New("帖子不存在")
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

	// Notify post author
	if post.AuthorID != authorID {
		s.notifSvc.Create(post.AuthorID, &authorID, model.NotifyComment,
			"有人评论了你的帖子《"+post.Title+"》", &input.PostID, "post")
	}

	// Notify parent comment author
	if input.ParentID != nil {
		parent, err := s.commentRepo.FindByID(*input.ParentID)
		if err == nil && parent.AuthorID != authorID {
			s.notifSvc.Create(parent.AuthorID, &authorID, model.NotifyReply,
				"有人回复了你的评论", input.ParentID, "comment")
		}
	}

	return s.commentRepo.FindByID(comment.ID)
}

func (s *CommentService) Delete(commentID, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	if comment.AuthorID != userID && !isAdmin {
		return errors.New("无权限删除此评论")
	}
	return s.commentRepo.Delete(commentID)
}

func (s *CommentService) List(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.ListByPost(postID, page, pageSize)
}

func (s *CommentService) GetCommentCount(postID uint) (int64, error) {
	return s.commentRepo.CountByPost(postID)
}
