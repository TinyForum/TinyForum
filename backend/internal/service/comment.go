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
	questionSvc *QuestionService // 添加问答服务依赖（可选）
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

	// 验证父评论
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

// ========== 新增的问答相关方法 ==========

// CreateAnswer 创建回答（仅限问答帖）
func (s *CommentService) CreateAnswer(authorID uint, input CreateCommentInput) (*model.Comment, error) {
	// Validate post exists and is question
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
		IsAnswer: true, // 标记为答案
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 回答奖励积分
	_ = s.userRepo.AddScore(authorID, 2)

	// Notify post author
	if post.AuthorID != authorID {
		s.notifSvc.Create(post.AuthorID, &authorID, model.NotifyComment,
			"有人回答了你的问题《"+post.Title+"》", &input.PostID, "post")
	}

	return s.commentRepo.FindByID(comment.ID)
}

// MarkAsAnswer 标记/取消标记为答案
func (s *CommentService) MarkAsAnswer(commentID, userID uint, isAdmin bool, isAnswer bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// Get post to check permissions
	post, err := s.postRepo.FindByID(comment.PostID)
	if err != nil {
		return errors.New("帖子不存在")
	}

	// Only post author or admin can mark as answer
	if post.AuthorID != userID && !isAdmin {
		return errors.New("无权限操作")
	}

	return s.commentRepo.MarkAsAnswer(commentID, isAnswer)
}

// GetAnswersByPostID 获取帖子的所有答案
func (s *CommentService) GetAnswersByPostID(postID uint, page, pageSize int, sortBy string) ([]model.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	switch sortBy {
	case "newest":
		return s.commentRepo.GetAnswersByPostIDOrderByNewest(postID, pageSize, offset)
	case "oldest":
		return s.commentRepo.GetAnswersByPostIDOrderByOldest(postID, pageSize, offset)
	default: // vote
		return s.commentRepo.GetAnswersByPostID(postID, pageSize, offset)
	}
}

// GetAnswerVoteCount 获取答案的投票数
func (s *CommentService) GetAnswerVoteCount(commentID uint) (int, error) {
	// 这个会通过 repository 的 questionRepo 来实现
	// 或者直接从 comment 的 vote_count 字段获取
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return 0, err
	}
	return comment.VoteCount, nil
}
