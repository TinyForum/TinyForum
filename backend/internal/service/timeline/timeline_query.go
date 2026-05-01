package timeline

import "tiny-forum/internal/model/po"

// GetHomeTimeline 获取首页时间线（推荐/综合）
func (s *timelineService) GetHomeTimeline(userID uint, page, pageSize int) ([]po.TimelineEvent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	_ = s.timelineRepo.UpdateLastRead(userID, "home")
	return s.timelineRepo.GetUserTimeline(userID, pageSize, offset)
}

// GetFollowingTimeline 获取关注时间线（仅关注用户的内容）
func (s *timelineService) GetFollowingTimeline(userID uint, page, pageSize int) ([]po.TimelineEvent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	_ = s.timelineRepo.UpdateLastRead(userID, "following")
	return s.timelineRepo.GetFollowingTimeline(userID, pageSize, offset)
}
