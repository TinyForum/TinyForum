package recommendation

import (
	"context"
	"time"

	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// RecommendationRepository 推荐系统数据访问接口
type RecommendationRepository interface {
	// 行为事件
	RecordBehavior(event *do.UserBehaviorEvent) error
	BatchRecordBehaviors(events []do.UserBehaviorEvent) error
	GetUserRecentBehaviors(userID uint, limit int) ([]do.UserBehaviorEvent, error)
	GetUserBehaviorsByType(userID uint, behaviorType do.BehaviorType, since time.Time) ([]do.UserBehaviorEvent, error)
	GetTargetBehaviors(targetID uint, targetType string) ([]do.UserBehaviorEvent, error)

	// 内容特征
	UpsertContentFeature(feature *do.ContentFeature) error
	GetContentFeature(creationID uint) (*do.ContentFeature, error)
	GetHotContentFeatures(limit int, excludeIDs []uint) ([]do.ContentFeature, error)
	GetContentFeaturesByBoard(boardID uint, limit int) ([]do.ContentFeature, error)
	GetContentFeaturesByTags(tagIDs []uint, excludeIDs []uint, limit int) ([]do.ContentFeature, error)

	// 推荐反馈
	RecordFeedback(fb *do.RecommendationFeedback) error
	BatchRecordFeedbacks(fbs []do.RecommendationFeedback) error
	GetUserFeedbacks(userID uint, since time.Time) ([]do.RecommendationFeedback, error)
	GetFeedbackCount(userID uint, feedbackType do.FeedbackType, since time.Time) (int64, error)

	// 用户兴趣画像
	UpsertInterestProfile(profile *do.UserInterestProfile) error
	GetInterestProfile(userID uint) (*do.UserInterestProfile, error)

	// 协同过滤：查找相似用户
	FindSimilarUsers(userID uint, limit int) ([]uint, error)
}

type recommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

// --- 行为事件 ---

func (r *recommendationRepository) RecordBehavior(event *do.UserBehaviorEvent) error {
	return r.db.Create(event).Error
}

func (r *recommendationRepository) BatchRecordBehaviors(events []do.UserBehaviorEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.CreateInBatches(events, 100).Error
}

func (r *recommendationRepository) GetUserRecentBehaviors(userID uint, limit int) ([]do.UserBehaviorEvent, error) {
	var events []do.UserBehaviorEvent
	err := r.db.Where("user_id = ?", userID).
		Order("created_ts DESC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *recommendationRepository) GetUserBehaviorsByType(userID uint, behaviorType do.BehaviorType, since time.Time) ([]do.UserBehaviorEvent, error) {
	var events []do.UserBehaviorEvent
	err := r.db.Where("user_id = ? AND behavior_type = ? AND created_at >= ?", userID, behaviorType, since).
		Order("created_ts DESC").
		Find(&events).Error
	return events, err
}

func (r *recommendationRepository) GetTargetBehaviors(targetID uint, targetType string) ([]do.UserBehaviorEvent, error) {
	var events []do.UserBehaviorEvent
	err := r.db.Where("target_id = ? AND target_type = ?", targetID, targetType).Find(&events).Error
	return events, err
}

// --- 内容特征 ---

func (r *recommendationRepository) UpsertContentFeature(feature *do.ContentFeature) error {
	return r.db.Save(feature).Error
}

func (r *recommendationRepository) GetContentFeature(creationID uint) (*do.ContentFeature, error) {
	var f do.ContentFeature
	err := r.db.Where("creation_id = ?", creationID).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *recommendationRepository) GetHotContentFeatures(limit int, excludeIDs []uint) ([]do.ContentFeature, error) {
	var features []do.ContentFeature
	q := r.db.Model(&do.ContentFeature{}).Order("hot_score DESC").Limit(limit)
	if len(excludeIDs) > 0 {
		q = q.Where("creation_id NOT IN ?", excludeIDs)
	}
	err := q.Find(&features).Error
	return features, err
}

func (r *recommendationRepository) GetContentFeaturesByBoard(boardID uint, limit int) ([]do.ContentFeature, error) {
	var features []do.ContentFeature
	err := r.db.Where("board_id = ?", boardID).Order("hot_score DESC").Limit(limit).Find(&features).Error
	return features, err
}

func (r *recommendationRepository) GetContentFeaturesByTags(tagIDs []uint, excludeIDs []uint, limit int) ([]do.ContentFeature, error) {
	if len(tagIDs) == 0 {
		return nil, nil
	}
	var features []do.ContentFeature
	tagJSONs := make([]string, len(tagIDs))
	for i, tid := range tagIDs {
		tagJSONs[i] = string(rune(tid))
	}
	q := r.db.Model(&do.ContentFeature{}).
		Where("tag_ids_json LIKE ?", "%"+tagJSONs[0]+"%").
		Order("hot_score DESC").Limit(limit)
	if len(excludeIDs) > 0 {
		q = q.Where("creation_id NOT IN ?", excludeIDs)
	}
	err := q.Find(&features).Error
	return features, err
}

// --- 推荐反馈 ---

func (r *recommendationRepository) RecordFeedback(fb *do.RecommendationFeedback) error {
	return r.db.Create(fb).Error
}

func (r *recommendationRepository) BatchRecordFeedbacks(fbs []do.RecommendationFeedback) error {
	if len(fbs) == 0 {
		return nil
	}
	return r.db.CreateInBatches(fbs, 100).Error
}

func (r *recommendationRepository) GetUserFeedbacks(userID uint, since time.Time) ([]do.RecommendationFeedback, error) {
	var fbs []do.RecommendationFeedback
	err := r.db.Where("user_id = ? AND created_at >= ?", userID, since).
		Order("created_at DESC").
		Find(&fbs).Error
	return fbs, err
}

func (r *recommendationRepository) GetFeedbackCount(userID uint, feedbackType do.FeedbackType, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&do.RecommendationFeedback{}).
		Where("user_id = ? AND feedback_type = ? AND created_at >= ?", userID, feedbackType, since).
		Count(&count).Error
	return count, err
}

// --- 用户兴趣画像 ---

func (r *recommendationRepository) UpsertInterestProfile(profile *do.UserInterestProfile) error {
	return r.db.Save(profile).Error
}

func (r *recommendationRepository) GetInterestProfile(userID uint) (*do.UserInterestProfile, error) {
	var p do.UserInterestProfile
	err := r.db.Where("user_id = ?", userID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// --- 协同过滤：查找与目标用户行为相似的用户 ---

func (r *recommendationRepository) FindSimilarUsers(userID uint, limit int) ([]uint, error) {
	var similarUserIDs []uint
	err := r.db.Raw(`
		SELECT DISTINCT ube2.user_id
		FROM user_behavior_events ube1
		INNER JOIN user_behavior_events ube2 ON ube1.target_id = ube2.target_id
		  AND ube1.target_type = ube2.target_type
		  AND ube1.behavior_type = ube2.behavior_type
		  AND ube2.user_id != ?
		WHERE ube1.user_id = ?
		  AND ube1.deleted_at IS NULL
		  AND ube2.deleted_at IS NULL
		GROUP BY ube2.user_id
		ORDER BY COUNT(DISTINCT ube1.target_id) DESC
		LIMIT ?
	`, userID, userID, limit).Scan(&similarUserIDs).Error
	return similarUserIDs, err
}

// Ensure RecommendationRepository is only used through the interface
var _ RecommendationRepository = (*recommendationRepository)(nil)

// Ensure context import is used (for future methods)
var _ context.Context
