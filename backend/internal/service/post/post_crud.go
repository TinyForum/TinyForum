package post

import (
	"errors"

	"tiny-forum/internal/dto"
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/model"

	"github.com/gin-gonic/gin"
)

type CreatePostInput struct {
	Title   string `json:"title" binding:"required,min=2,max=200"`
	Content string `json:"content" binding:"required,min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	Type    string `json:"type"`
	BoardID uint   `json:"board_id" binding:"required"`
	TagIDs  []uint `json:"tag_ids"`
}

type UpdatePostInput struct {
	Title   string `json:"title" binding:"min=2,max=200"`
	Content string `json:"content" binding:"min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	TagIDs  []uint `json:"tag_ids"`
}

// Create 创建帖子
func (s *PostService) Create(ctx *gin.Context, authorID uint, input CreatePostInput) (*model.Post, error) {
	// postType := model.PostType(input.Type)
	// 默认创建帖子
	postType := model.PostType(input.Type)
	if postType == "" || !postType.IsValid() {
		postType = model.PostTypePost
	}

	// 发布板块检查
	if input.BoardID == 0 {
		return nil, errors.New("board_id 不能为 0")
	}
	board, err := s.boardRepo.FindByID(input.BoardID)
	if err != nil {
		return nil, errors.New("选择的板块不存在")
	}

	// 是否需要审核
	reviewRequired, hitWords := middleware.IsReviewRequired(ctx)
	// 是否屏蔽
	shadowed, hitWords := middleware.IsShadowed(ctx)

	replace, hitWords := middleware.IsReplaced(ctx)

	moderationStatus := model.AuditStatusPending
	if reviewRequired {
		// 先标记为待审核状态
		moderationStatus = model.AuditStatusPending
	}
	if shadowed {
		moderationStatus = model.AuditStatusRejected
	}
	if replace {
		moderationStatus = model.AuditStatusApproved
	}

	post := &model.Post{
		Title:   input.Title,
		Content: input.Content,
		Summary: input.Summary,
		Cover:   input.Cover,
		Type:    postType,
		// Status:   model.PostStatusPublished,
		AuthorID: authorID,
		BoardID:  board.ID,
		// PostStatus:       poststatus,
		ModerationStatus: moderationStatus,
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
	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}
	for _, tag := range post.Tags {
		_ = s.tagRepo.IncrPostCount(tag.ID, 1)
	}
	_ = s.userRepo.AddScore(authorID, 10)
	if post, err = s.postRepo.FindByID(post.ID); err != nil {
		return nil, err
	}

	if reviewRequired || shadowed || replace {
		go func() {
			_ = s.contentcheckSvc.CreateAuditTaskForPost(post.ID, "sensitive_word", hitWords)
		}()
	}
	return post, nil
}

// Update 更新帖子
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

// Delete 删除帖子
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

// GetByID 获取帖子详情（含点赞状态）
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

// List 获取帖子列表（支持筛选）
func (s *PostService) List(page, pageSize int, opts dto.PostListOptions) ([]model.Post, int64, error) {
	return s.postRepo.List(page, pageSize, opts)
}
