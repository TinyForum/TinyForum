package board

import (
	"fmt"
	"regexp"

	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service/notification"
)

type BoardService struct {
	boardRepo *repository.BoardRepository
	userRepo  *repository.UserRepository
	postRepo  repository.PostRepository
	notifSvc  *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewBoardService(
	boardRepo *repository.BoardRepository,
	userRepo *repository.UserRepository,
	postRepo repository.PostRepository,
	notifSvc *notification.NotificationService,
) *BoardService {
	return &BoardService{
		boardRepo: boardRepo,
		userRepo:  userRepo,
		postRepo:  postRepo,
		notifSvc:  notifSvc,
	}
}

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
	valid := map[model.UserRole]bool{
		model.RoleGuest:     true,
		model.RoleUser:      true,
		model.RoleMember:    true,
		model.RoleModerator: true,
		model.RoleAdmin:     true,
	}
	for _, r := range roles {
		if r != "" && !valid[model.UserRole(r)] {
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
func (s *BoardService) writeLog(moderatorID, boardID uint, action, targetType string, targetID uint, reason string) {
	log := &model.ModeratorLog{
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
func (s *BoardService) writeLogWithValues(moderatorID, boardID uint, action, targetType string, targetID uint, reason, oldValue, newValue string) {
	log := &model.ModeratorLog{
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
