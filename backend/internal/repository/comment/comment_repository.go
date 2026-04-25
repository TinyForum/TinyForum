package comment

import (
	"errors"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// Create 创建评论
func (r *commentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID 根据 ID 获取评论（预加载作者）
func (r *commentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("Author").First(&comment, id).Error
	return &comment, err
}

// Update 更新评论
func (r *commentRepository) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// ValidateParentComment 验证父评论是否属于同一帖子
func (r *commentRepository) ValidateParentComment(parentID uint, postID uint) error {
	var comment model.Comment
	err := r.db.Where("id = ? AND post_id = ?", parentID, postID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("父评论不存在或不属于当前帖子")
		}
		return err
	}
	return nil
}
