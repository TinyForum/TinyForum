package tag

import (
	"tiny-forum/internal/model"
)

// FindTagsByPostIDs 批量查询帖子的标签
func (r *tagRepository) FindTagsByPostIDs(postIDs []uint) (map[uint][]model.Tag, error) {
	if len(postIDs) == 0 {
		return make(map[uint][]model.Tag), nil
	}

	type PostTagRelation struct {
		PostID uint
		Tag    model.Tag
	}

	var relations []PostTagRelation
	err := r.db.Table("post_tags").
		Select("post_tags.post_id, tags.*").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id IN ?", postIDs).
		Scan(&relations).Error
	if err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]model.Tag)
	for _, rel := range relations {
		tagMap[rel.PostID] = append(tagMap[rel.PostID], rel.Tag)
	}
	return tagMap, nil
}

// FindTagsByPostID 查询单个帖子的标签
func (r *tagRepository) FindTagsByPostID(postID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Table("tags").
		Select("tags.*").
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error
	return tags, err
}
