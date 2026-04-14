package service

import (
	"errors"
	"fmt"

	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo  repository.PostRepository
	tagRepo   *repository.TagRepository
	boardRepo *repository.BoardRepository
	userRepo  *repository.UserRepository
	notifSvc  *NotificationService
}

func NewPostService(
	postRepo repository.PostRepository,
	tagRepo *repository.TagRepository,
	userRepo *repository.UserRepository,
	boardRepo *repository.BoardRepository,
	notifSvc *NotificationService,
) *PostService {
	return &PostService{postRepo: postRepo, tagRepo: tagRepo, userRepo: userRepo, boardRepo: boardRepo, notifSvc: notifSvc}
}

type CreatePostInput struct {
	Title   string `json:"title" binding:"required,min=2,max=200"`
	Content string `json:"content" binding:"required,min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	Type    string `json:"type"`
	BoardID uint   `json:"board_id" binding:"required"`
	TagIDs  []uint `json:"tag_ids"`
}

func (s *PostService) Create(authorID uint, input CreatePostInput) (*model.Post, error) {
	postType := model.PostType(input.Type)
	if postType == "" {
		postType = model.PostTypePost
	}
	// 验证 board_id
	if input.BoardID == 0 {
		return nil, errors.New("board_id 不能为 0")
	}

	// 验证板块是否存在
	board, err := s.boardRepo.FindByID(input.BoardID)
	if err != nil {
		return nil, errors.New("选择的板块不存在")
	}

	// 验证帖子类型
	validTypes := map[string]bool{"post": true, "article": true, "topic": true}
	if !validTypes[input.Type] {
		input.Type = "post"
	}

	post := &model.Post{
		Title:    input.Title,
		Content:  input.Content,
		Summary:  input.Summary,
		Cover:    input.Cover,
		Type:     postType,
		Status:   model.PostStatusPublished,
		AuthorID: authorID,
		BoardID:  board.ID, // 使用从数据库查询到的板块 ID
	}

	// Attach tags
	if len(input.TagIDs) > 0 {
		var tags []model.Tag
		for _, id := range input.TagIDs {
			tag, err := s.tagRepo.FindByID(id)
			if err == nil {
				tags = append(tags, *tag)
			}
		}
		post.Tags = tags
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// Increment tag post count
	for _, tag := range post.Tags {
		_ = s.tagRepo.IncrPostCount(tag.ID, 1)
	}

	// Score for author
	_ = s.userRepo.AddScore(authorID, 10)

	return s.postRepo.FindByID(post.ID)
}

type UpdatePostInput struct {
	Title   string `json:"title" binding:"min=2,max=200"`
	Content string `json:"content" binding:"min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	TagIDs  []uint `json:"tag_ids"`
}

func (s *PostService) Update(postID, userID uint, isAdmin bool, input UpdatePostInput) (*model.Post, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}
	if post.AuthorID != userID && !isAdmin {
		return nil, errors.New("无权限修改此帖子")
	}

	if input.Title != "" {
		post.Title = input.Title
	}
	if input.Content != "" {
		post.Content = input.Content
	}
	if input.Summary != "" {
		post.Summary = input.Summary
	}
	if input.Cover != "" {
		post.Cover = input.Cover
	}

	if len(input.TagIDs) > 0 {
		var tags []model.Tag
		for _, id := range input.TagIDs {
			tag, err := s.tagRepo.FindByID(id)
			if err == nil {
				tags = append(tags, *tag)
			}
		}
		post.Tags = tags
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}
	return s.postRepo.FindByID(post.ID)
}

func (s *PostService) Delete(postID, userID uint, isAdmin bool) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.AuthorID != userID && !isAdmin {
		return errors.New("无权限删除此帖子")
	}
	return s.postRepo.Delete(postID)
}

func (s *PostService) GetByID(postID, viewerID uint) (*model.Post, bool, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, false, errors.New("帖子不存在")
	}
	_ = s.postRepo.IncrViewCount(postID)

	liked := false
	if viewerID > 0 {
		liked = s.postRepo.IsLiked(viewerID, postID)
	}
	return post, liked, nil
}

func (s *PostService) List(page, pageSize int, opts repository.PostListOptions) ([]model.Post, int64, error) {
	return s.postRepo.List(page, pageSize, opts)
}

func (s *PostService) Like(userID, postID uint) error {
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

func (s *PostService) Unlike(userID, postID uint) error {
	if err := s.postRepo.RemoveLike(userID, postID); err != nil {
		return err
	}
	return s.postRepo.IncrLikeCount(postID, -1)
}

func (s *PostService) AdminList(page, pageSize int, keyword string) ([]model.Post, int64, error) {
	return s.postRepo.AdminList(page, pageSize, keyword)
}

func (s *PostService) SetStatus(postID uint, status model.PostStatus) error {
	return s.postRepo.Update(&model.Post{})
}

func (s *PostService) TogglePin(postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrPostNotFound // 返回标准错误
		}
		return fmt.Errorf("查询帖子失败: %w", err)
	}
	post.PinTop = !post.PinTop
	return s.postRepo.Update(post)
}
