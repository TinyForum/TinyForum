package article

import (
	"context"
	"errors"
	"fmt"

	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/converter"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// AdminList 管理员获取帖子列表
func (s *articleService) AdminLists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListPosts]) ([]do.Article, int64, error) {
	// 将 BO 的 Data 字段转换为 DO 的对应结构
	// var filterDO *do.Post
	filterDO := converter.ListPostsBOToPostDO(&listPostsBO.Data)

	// 构造 DO 层的查询对象
	listPostsDO := &common.PageQuery[do.Article]{
		Page:     listPostsBO.Page,
		PageSize: listPostsBO.PageSize,
		Data:     *filterDO,
		Keyword:  listPostsBO.Keyword,
		SortBy:   listPostsBO.SortBy,
		TagNames: listPostsBO.TagNames,
	}

	// 如果 filterDO 可能为 nil，且 Data 字段要求非指针，则需要处理零值
	if filterDO == nil {
		listPostsDO.Data = do.Article{}
	} else {
		listPostsDO.Data = *filterDO
	}
	filterDO.PostStatus = do.PostStatusPublished

	return s.postRepo.AdminList(ctx, listPostsDO)
}

// SetStatus 设置帖子状态（暂未完全实现，保留接口）
func (s *articleService) SetStatus(postID uint, status do.PostStatus) error {
	return s.postRepo.Update(&do.Article{})
}

// TogglePin 切换帖子置顶状态
func (s *articleService) TogglePin(postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrPostNotFound
		}
		return fmt.Errorf("查询帖子失败: %w", err)
	}
	post.PinTop = !post.PinTop
	return s.postRepo.Update(post)
}

// 管理员更新审核状态
func (s *articleService) AdminSetReviewPost(postID uint, status do.ModerationStatus) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return gorm.ErrRecordNotFound
	}

	post.ModerationStatus = do.ModerationStatus(status)
	return s.postRepo.Update(post)
}
