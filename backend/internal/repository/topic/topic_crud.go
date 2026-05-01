package topic

import "tiny-forum/internal/model/do"

func (r *topicRepository) Create(topic *do.Topic) error {
	return r.db.Create(topic).Error
}

func (r *topicRepository) Update(topic *do.Topic) error {
	return r.db.Save(topic).Error
}

func (r *topicRepository) Delete(id uint) error {
	return r.db.Delete(&do.Topic{}, id).Error
}

func (r *topicRepository) FindByID(id uint) (*do.Topic, error) {
	var topic do.Topic
	err := r.db.Preload("Creator").First(&topic, id).Error
	return &topic, err
}

func (r *topicRepository) List(limit, offset int) ([]do.Topic, int64, error) {
	var topics []do.Topic
	var total int64

	query := r.db.Model(&do.Topic{})
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("follower_count DESC, post_count DESC").
		Find(&topics).Error
	return topics, total, err
}

func (r *topicRepository) GetByCreator(creatorID uint, limit, offset int) ([]do.Topic, int64, error) {
	var topics []do.Topic
	var total int64

	query := r.db.Model(&do.Topic{}).Where("creator_id = ?", creatorID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&topics).Error
	return topics, total, err
}
