package comment

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// Create 创建评论
func (r *commentRepository) CreateComment(comment *do.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) CreateAnswer(comment *do.Answer) error {
	return r.db.Create(comment).Error
}

// FindByCommentID 根据 ID 获取评论（预加载作者）
func (r *commentRepository) FindByCommentID(id uint) (*do.Comment, error) {
	var comment do.Comment
	err := r.db.Preload("Reply.Author").First(&comment, id).Error
	return &comment, err
}

// 根据 ID 获取回答（预加载作者）
func (r *commentRepository) FindByAnswerID(id uint) (*do.Answer, error) {
	var answer do.Answer
	err := r.db.Preload("Reply.Author").First(&answer, id).Error
	return &answer, err
}

// Update 更新评论
func (r *commentRepository) Update(comment *do.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&do.Comment{}, id).Error
}

// ValidateParentComment 验证父评论是否属于同一帖子
func (r *commentRepository) ValidateParentComment(parentID uint, postID uint) error {
	var reply do.Reply // 注意替换为 Reply 模型，而不是 Comment
	err := r.db.Where("id = ? AND target_id = ?", parentID, postID).First(&reply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("父评论不存在或不属于当前帖子")
		}
		return err
	}
	return nil
}

// GetCommentTree 获取指定帖子下的完整评论树（递归加载所有子评论）
func (r *commentRepository) GetCommentTree(postID uint) ([]do.Comment, error) {
	var comments []do.Comment

	// 1. 获取顶层评论：works_id = postID 且关联的 Reply 的 parent_id IS NULL
	err := r.db.
		Preload("Reply.Author"). // 预加载作者（顶层）
		Joins("JOIN replies ON replies.id = comments.reply_id").
		Where("comments.works_id = ? AND replies.parent_id IS NULL", postID).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// 2. 对每个顶层评论的 Reply，递归加载其所有后代 Reply
	for i := range comments {
		if comments[i].Reply == nil {
			continue // 数据异常，跳过
		}
		if err := r.loadRepliesRecursive(comments[i].Reply); err != nil {
			return nil, err
		}
	}
	return comments, nil
}

// loadRepliesRecursive 递归加载一个 Reply 的所有后代 Reply（预加载作者）
func (r *commentRepository) loadRepliesRecursive(reply *do.Reply) error {
	var children []do.Reply
	err := r.db.
		Preload("Author"). // 预加载子回复的作者
		Where("parent_id = ?", reply.ID).
		Find(&children).Error
	if err != nil {
		return err
	}
	reply.Replies = children

	// 递归加载每个子回复的后代
	for i := range children {
		if err := r.loadRepliesRecursive(&children[i]); err != nil {
			return err
		}
	}
	return nil
}
