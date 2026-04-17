package post

import (
	"errors"
	"fmt"

	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// AdminList 管理员获取帖子列表
func (s *PostService) AdminList(page, pageSize int, keyword string) ([]model.Post, int64, error) {
	return s.postRepo.AdminList(page, pageSize, keyword)
}

// SetStatus 设置帖子状态（暂未完全实现，保留接口）
func (s *PostService) SetStatus(postID uint, status model.PostStatus) error {
	return s.postRepo.Update(&model.Post{})
}

// TogglePin 切换帖子置顶状态
func (s *PostService) TogglePin(postID uint) error {
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
