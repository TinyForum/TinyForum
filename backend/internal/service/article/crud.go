package article

import (
	"context"
	"errors"

	"tiny-forum/internal/middleware"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/converter"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Create 创建帖子
func (s *articleService) Create(ctx *gin.Context, authorID uint, input request.CreatePostRequest) (*do.Article, error) {
	// 1. 帖子类型校验
	postType := do.CreationType(input.Type)
	if postType == "" || !postType.IsValid() {
		postType = do.CreationTypePost
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
	var moderationStatus do.ModerationStatus
	switch {
	case shadowed:
		moderationStatus = do.ModerationStatusRejected
	case reviewRequired:
		moderationStatus = do.ModerationStatusPending
	case replaced:
		moderationStatus = do.ModerationStatusApproved
	default:
		moderationStatus = do.ModerationStatusApproved
	}

	// 6. 构建帖子对象
	post := &do.Article{
		Creation: do.Creation{
			Title:            input.Title,
			Content:          input.Content,
			Summary:          input.Summary,
			CoverUrl:         input.Cover,
			Slug:             utils.GenerateSlug(),
			Type:             postType,
			AuthorID:         authorID,
			BoardID:          board.ID,
			ModerationStatus: moderationStatus,
			CreationStatus:   do.CreationStatus(input.Status),
		},
	}

	// 7. 处理标签
	if len(input.TagIDs) > 0 {
		tags := make([]do.Tag, 0, len(input.TagIDs))
		for _, id := range input.TagIDs {
			tag, err := s.tagRepo.FindByID(id)
			if err == nil {
				tags = append(tags, *tag)
			}
		}
		post.Creation.Tags = tags
	}

	// 8. 创建帖子
	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// 9. 更新标签计数
	for _, tag := range post.Creation.Tags {
		_ = s.tagRepo.IncrPostCount(tag.ID, 1)
	}

	// 10. 增加用户积分
	_ = s.userRepo.AddScore(authorID, 10)

	// 11. 重新加载完整帖子（包含关联数据）
	post, err = s.postRepo.FindByArticleID(post.ID)
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
func (s *articleService) Update(postID, userID uint, isAdmin bool, input request.UpdatePostRequest) (*do.Article, error) {
	post, err := s.postRepo.FindByArticleID(postID)
	if err != nil {
		return nil, apperrors.ErrPostNotFound
	}
	if post.Creation.AuthorID != userID && !isAdmin {
		return nil, apperrors.ErrInsufficientPermission
	}
	if input.Title != "" {
		post.Creation.Title = input.Title
	}
	if input.Content != "" {
		post.Creation.Content = input.Content
	}
	if input.Summary != "" {
		post.Creation.Summary = input.Summary
	}
	if input.Cover != "" {
		post.Creation.CoverUrl = input.Cover
	}
	if len(input.TagIDs) > 0 {
		var tags []do.Tag
		for _, id := range input.TagIDs {
			tag, err := s.tagRepo.FindByID(id)
			if err == nil {
				tags = append(tags, *tag)
			}
		}
		post.Creation.Tags = tags
	}
	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}
	return s.postRepo.FindByArticleID(post.ID)
}

// Delete 删除帖子
func (s *articleService) Delete(postID, userID uint, isAdmin bool) error {
	post, err := s.postRepo.FindByArticleID(postID)
	if err != nil {
		return apperrors.ErrPostNotFound
	}
	if post.Creation.AuthorID != userID && !isAdmin {
		return apperrors.ErrInsufficientPermission
	}
	return s.postRepo.Delete(postID)
}

// GetByID 获取帖子详情（含点赞状态）
func (s *articleService) GetByID(postID, viewerID uint) (*do.Article, bool, error) {
	post, err := s.postRepo.FindByArticleID(postID)
	if err != nil {
		return nil, false, apperrors.ErrPostNotFound
	}
	_ = s.postRepo.IncrViewCount(postID)
	liked := false
	if viewerID > 0 {
		liked = s.postRepo.IsLiked(viewerID, postID)
	}
	return post, liked, nil
}

// 用户获取文章列表
func (s *articleService) List(ctx context.Context, listPostsBO *common.PageQuery[bo.ListPosts]) ([]do.Article, int64, error) {
	filterDO := converter.ListPostsBOToPostDO(&listPostsBO.Data)
	// 假定转换器总是返回非空指针（若传入 nil 则返回 &do.Post{}）

	// 构造 DO 层的查询对象，外层查询参数直接赋值
	listPostsDO := &common.PageQuery[do.Article]{
		Page:     listPostsBO.Page,
		PageSize: listPostsBO.PageSize,
		Data:     *filterDO,
		Keyword:  listPostsBO.Keyword,
		SortBy:   listPostsBO.SortBy,
		TagNames: listPostsBO.TagNames,
	}

	return s.postRepo.List(ctx, listPostsDO)
}
