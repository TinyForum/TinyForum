package board

import (
	"fmt"
	"regexp"
	"tiny-forum/internal/model/po"
)

// validateSlug 校验 slug 格式
func validateSlug(slug string) error {
	if len(slug) == 0 || len(slug) > 50 {
		return fmt.Errorf("板块标识长度必须在1-50字符之间")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9\-_]+$`, slug)
	if !matched {
		return fmt.Errorf("板块标识只能包含小写字母、数字、横线和下划线")
	}
	return nil
}

// validateRoles 校验 ViewRole / PostRole / ReplyRole 是否合法
func validateRoles(roles ...string) error {
	valid := map[po.UserRole]bool{
		po.RoleGuest:     true,
		po.RoleUser:      true,
		po.RoleMember:    true,
		po.RoleModerator: true,
		po.RoleAdmin:     true,
	}
	for _, r := range roles {
		if r != "" && !valid[po.UserRole(r)] {
			return fmt.Errorf("无效的角色值: %s", r)
		}
	}
	return nil
}

// boolVal 返回指针值或 fallback
func boolVal(ptr *bool, fallback bool) bool {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

// writeLog 记录版主操作日志（忽略错误）
func (s *boardService) writeLog(moderatorID, boardID uint, action, targetType string, targetID uint, reason string) {
	log := &po.ModeratorLog{
		ModeratorID: moderatorID,
		BoardID:     boardID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
	}
	_ = s.boardRepo.CreateModeratorLog(log)
}

// writeLogWithValues 同 writeLog，额外写入 OldValue / NewValue
func (s *boardService) writeLogWithValues(moderatorID, boardID uint, action, targetType string, targetID uint, reason, oldValue, newValue string) {
	log := &po.ModeratorLog{
		ModeratorID: moderatorID,
		BoardID:     boardID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		OldValue:    oldValue,
		NewValue:    newValue,
	}
	_ = s.boardRepo.CreateModeratorLog(log)
}
