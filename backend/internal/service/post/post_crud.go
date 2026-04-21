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
	Status  string `json:"status"`
}

type UpdatePostInput struct {
	Title   string `json:"title" binding:"min=2,max=200"`
	Content string `json:"content" binding:"min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	TagIDs  []uint `json:"tag_ids"`
}

// Create 创建帖子
// Create 创建帖子
func (s *PostService) Create(ctx *gin.Context, authorID uint, input CreatePostInput) (*model.Post, error) {
	// 1. 帖子类型校验
	postType := model.PostType(input.Type)
	if postType == "" || !postType.IsValid() {
		postType = model.PostTypePost
	}

	// 2. 板块校验
	if input.BoardID == 0 {
		return nil, errors.New("board_id 不能为 0")
	}
	board, err := s.boardRepo.FindByID(input.BoardID)
	if err != nil {
		return nil, errors.New("选择的板块不存在")
	}

	// 3. 获取中间件注入的审核标记（分别获取，避免覆盖）
	reviewRequired, reviewHitWords := middleware.IsReviewRequired(ctx)
	shadowed, shadowHitWords := middleware.IsShadowed(ctx)
	replaced, replaceHitWords := middleware.IsReplaced(ctx)

	// 4. 合并所有命中词
	allHitWords := make([]string, 0,
		len(reviewHitWords)+len(shadowHitWords)+len(replaceHitWords))
	allHitWords = append(allHitWords, reviewHitWords...)
	allHitWords = append(allHitWords, shadowHitWords...)
	allHitWords = append(allHitWords, replaceHitWords...)

	// 5. 确定审核状态（优先级：屏蔽 > 待审 > 替换 > 安全）
	var moderationStatus model.ModerationStatus
	switch {
	case shadowed:
		moderationStatus = model.ModerationStatusRejected
	case reviewRequired:
		moderationStatus = model.ModerationStatusPending
	case replaced:
		moderationStatus = model.ModerationStatusApproved
	default:
		moderationStatus = model.ModerationStatusApproved
	}

	// 6. 构建帖子对象
	post := &model.Post{
		Title:            input.Title,
		Content:          input.Content,
		Summary:          input.Summary,
		Cover:            input.Cover,
		Type:             postType,
		AuthorID:         authorID,
		BoardID:          board.ID,
		ModerationStatus: moderationStatus,
		PostStatus:       model.PostStatus(input.Status),
	}

	// 7. 处理标签
	if len(input.TagIDs) > 0 {
		tags := make([]model.Tag, 0, len(input.TagIDs))
		for _, id := range input.TagIDs {
			tag, err := s.tagRepo.FindByID(id)
			if err == nil {
				tags = append(tags, *tag)
			}
		}
		post.Tags = tags
	}

	// 8. 创建帖子
	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// 9. 更新标签计数
	for _, tag := range post.Tags {
		_ = s.tagRepo.IncrPostCount(tag.ID, 1)
	}

	// 10. 增加用户积分
	_ = s.userRepo.AddScore(authorID, 10)

	// 11. 重新加载完整帖子（包含关联数据）
	post, err = s.postRepo.FindByID(post.ID)
	if err != nil {
		return nil, err
	}

	// 12. 异步创建审核任务（如有需要）
	if reviewRequired || shadowed || replaced {
		go func() {
			_ = s.contentcheckSvc.CreateAuditTaskForPost(post.ID, "sensitive_word", allHitWords)
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
