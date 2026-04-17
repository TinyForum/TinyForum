package topic

import (
	"tiny-forum/internal/model"
)

func (r *TopicRepository) Create(topic *model.Topic) error {
	return r.db.Create(topic).Error
}

func (r *TopicRepository) Update(topic *model.Topic) error {
	return r.db.Save(topic).Error
}

func (r *TopicRepository) Delete(id uint) error {
	return r.db.Delete(&model.Topic{}, id).Error
}

func (r *TopicRepository) FindByID(id uint) (*model.Topic, error) {
	var topic model.Topic
	err := r.db.Preload("Creator").First(&topic, id).Error
	return &topic, err
}

func (r *TopicRepository) List(limit, offset int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64

	query := r.db.Model(&model.Topic{})
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("follower_count DESC, post_count DESC").
		Find(&topics).Error
	return topics, total, err
}

func (r *TopicRepository) GetByCreator(creatorID uint, limit, offset int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64

	query := r.db.Model(&model.Topic{}).Where("creator_id = ?", creatorID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&topics).Error
	return topics, total, err
}
