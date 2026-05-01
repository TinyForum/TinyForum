package topic

import (
	"tiny-forum/internal/model/po"
)

func (r *topicRepository) Create(topic *po.Topic) error {
	return r.db.Create(topic).Error
}

func (r *topicRepository) Update(topic *po.Topic) error {
	return r.db.Save(topic).Error
}

func (r *topicRepository) Delete(id uint) error {
	return r.db.Delete(&po.Topic{}, id).Error
}

func (r *topicRepository) FindByID(id uint) (*po.Topic, error) {
	var topic po.Topic
	err := r.db.Preload("Creator").First(&topic, id).Error
	return &topic, err
}

func (r *topicRepository) List(limit, offset int) ([]po.Topic, int64, error) {
	var topics []po.Topic
	var total int64

	query := r.db.Model(&po.Topic{})
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("follower_count DESC, post_count DESC").
		Find(&topics).Error
	return topics, total, err
}

func (r *topicRepository) GetByCreator(creatorID uint, limit, offset int) ([]po.Topic, int64, error) {
	var topics []po.Topic
	var total int64

	query := r.db.Model(&po.Topic{}).Where("creator_id = ?", creatorID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&topics).Error
	return topics, total, err
}
