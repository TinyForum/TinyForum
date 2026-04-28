package notification

import (
	"errors"
	"tiny-forum/internal/bo"
	"tiny-forum/internal/model"
)

// Create 创建通知（忽略错误，不返回 error 保持原有行为）
func (s *notificationService) Create(userID uint, senderID *uint, notifType model.NotificationType, content string, targetID *uint, targetType string) {
	n := &model.Notification{
		UserID:     userID,
		SenderID:   senderID,
		Type:       notifType,
		Content:    content,
		TargetID:   targetID,
		TargetType: targetType,
	}
	_ = s.notifRepo.Create(n)
}

// List 获取用户通知列表（分页）
// List 获取通知列表（返回 BO）
func (s *notificationService) List(userID uint, page, pageSize int) (*bo.NotificationListResult, error) {
	notifications, total, err := s.notifRepo.ListByUser(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.notifRepo.UnreadCount(userID)
	if err != nil {
		return nil, err
	}

	return &bo.NotificationListResult{
		List:        modelsToBOs(notifications), // 使用转换函数
		Total:       total,
		UnreadCount: unreadCount,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// MarkAllRead 标记所有通知为已读
func (s *notificationService) MarkAllRead(userID uint) (int64, error) {
	return s.notifRepo.MarkAllRead(userID)
}

// UnreadCount 获取未读通知数量
func (s *notificationService) UnreadCount(userID uint) (int64, error) {
	return s.notifRepo.UnreadCount(userID)
}

func (s *notificationService) MarkRead(notificationID uint, userID uint) error {
	// 1. 查询通知（通过 Repository）
	notif, err := s.notifRepo.GetByID(notificationID)
	if err != nil {
		return errors.New("通知不存在")
	}

	// 2. 权限验证
	if notif.UserID != userID {
		return errors.New("无权操作此通知")
	}

	// 3. 标记已读
	return s.notifRepo.MarkRead(notificationID)
}

// BatchMarkRead 批量标记通知为已读
func (s *notificationService) BatchMarkRead(userID uint, ids []uint) (int64, error) {
	if len(ids) == 0 {
		// 标记所有为已读
		return s.notifRepo.MarkAllRead(userID)
	}
	// 批量标记指定通知
	return s.notifRepo.BatchMarkRead(userID, ids)
}
