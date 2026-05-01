package post

import (
	"errors"
	"fmt"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// AdminList 管理员获取帖子列表
func (s *postService) AdminList(page, pageSize int, opts dto.PostListOptions) ([]do.Post, int64, error) {
	return s.postRepo.AdminList(page, pageSize, opts)
}

// SetStatus 设置帖子状态（暂未完全实现，保留接口）
func (s *postService) SetStatus(postID uint, status do.PostStatus) error {
	return s.postRepo.Update(&do.Post{})
}

// TogglePin 切换帖子置顶状态
func (s *postService) TogglePin(postID uint) error {
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
func (s *postService) AdminSetReviewPost(postID uint, status do.ModerationStatus) error {
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
