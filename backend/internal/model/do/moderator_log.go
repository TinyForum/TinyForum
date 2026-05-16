package do

import (
	"fmt"
	"tiny-forum/internal/model/common"
)

// ActionType 版主操作类型
type ModeratorActionType string

const (
	ActionDeletePost     ModeratorActionType = "delete_post"      // 删除帖子
	ActionDeleteComment  ModeratorActionType = "delete_comment"   // 删除评论
	ActionBanUser        ModeratorActionType = "ban_user"         // 封禁用户
	ActionUnbanUser      ModeratorActionType = "unban_user"       // 解封用户
	ActionPinPost        ModeratorActionType = "pin_post"         // 置顶帖子
	ActionUnpinPost      ModeratorActionType = "unpin_post"       // 取消置顶
	ActionLockPost       ModeratorActionType = "lock_post"        // 锁定帖子
	ActionUnlockPost     ModeratorActionType = "unlock_post"      // 解锁帖子
	ActionSetBoardConfig ModeratorActionType = "set_board_config" // 修改版块配置
)

// IsValid 验证操作类型是否有效
func (a ModeratorActionType) IsValid() bool {
	switch a {
	case ActionDeletePost, ActionDeleteComment, ActionBanUser, ActionUnbanUser,
		ActionPinPost, ActionUnpinPost, ActionLockPost, ActionUnlockPost, ActionSetBoardConfig:
		return true
	}
	return false
}

// ParseActionType 严格解析操作类型，无效时返回错误
func ParsemoderatorActionType(s string) (ModeratorActionType, error) {
	at := ModeratorActionType(s)
	if at.IsValid() {
		return at, nil
	}
	return "", fmt.Errorf("invalid action type: %s", s)
}

// TargetType 版主操作的目标类型
type ModeratorTargetType string

const (
	TargetPost    ModeratorTargetType = "post"    // 帖子
	TargetComment ModeratorTargetType = "comment" // 评论
	TargetUser    ModeratorTargetType = "user"    // 用户
	TargetBoard   ModeratorTargetType = "board"   // 版块
)

// IsValid 验证目标类型是否有效
func (tt ModeratorTargetType) IsValid() bool {
	switch tt {
	case TargetPost, TargetComment, TargetUser, TargetBoard:
		return true
	}
	return false
}

// ParseTargetType 严格解析目标类型，无效时返回错误
func ParseModeratorTargetType(s string) (ModeratorTargetType, error) {
	tt := ModeratorTargetType(s)
	if tt.IsValid() {
		return tt, nil
	}
	return "", fmt.Errorf("invalid target type: %s", s)
}

// ModeratorLog 版主操作日志
type ModeratorLog struct {
	common.BaseModel
	ModeratorID uint                `gorm:"not null;index:idx_moderator_time,priority:1" json:"moderator_id"` // 操作者 ID
	BoardID     uint                `gorm:"index:idx_board_time,priority:1" json:"board_id"`                  // 版块 ID
	Action      ModeratorActionType `gorm:"type:varchar(50);not null;index" json:"action"`                    // 操作类型
	TargetType  ModeratorTargetType `gorm:"type:varchar(20);not null" json:"target_type"`                     // 目标类型
	TargetID    uint                `gorm:"not null;index" json:"target_id"`                                  // 目标 ID
	Reason      string              `gorm:"type:text" json:"reason"`                                          // 操作理由，支持长文本
	OldValue    string              `gorm:"type:json" json:"old_value"`                                       // 变更前值（JSON）
	NewValue    string              `gorm:"type:json" json:"new_value"`                                       // 变更后值（JSON）

	// 关联字段：使用指针避免 JSON 空对象，omitempty 在指针为 nil 时生效
	Moderator *User  `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
	Board     *Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}

func (ModeratorLog) TableName() string { return "moderator_logs" }
