package repository

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func (r *TagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("post_count DESC").Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

func (r *TagRepository) IncrPostCount(id uint, delta int) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}

// internal/repository/tag_repository.go

// FindTagsByPostIDs 批量查询帖子的标签
func (r *TagRepository) FindTagsByPostIDs(postIDs []uint) (map[uint][]model.Tag, error) {
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

	// 组装为 map
	tagMap := make(map[uint][]model.Tag)
	for _, rel := range relations {
		tagMap[rel.PostID] = append(tagMap[rel.PostID], rel.Tag)
	}
	return tagMap, nil
}

// FindTagsByPostID 查询单个帖子的标签
func (r *TagRepository) FindTagsByPostID(postID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Table("tags").
		Select("tags.*").
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error
	return tags, err
}

// internal/repository/tag_repository.go

// Count 获取标签总数
func (r *TagRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).Count(&count).Error
	return count, err
}

// CountByDateRange 根据日期范围统计标签数
func (r *TagRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}
