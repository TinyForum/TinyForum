package request

import (
	"fmt"
	"tiny-forum/internal/model/do"

	"gorm.io/datatypes"
)

type ApplyModeratorRequest struct {
	UserID               uint                                        `json:"user_id"`               // 申请人ID
	BoardID              uint                                        `json:"board_id"`              // 版块ID
	Reason               string                                      `json:"reason"`                // 申请理由
	RequestedPermissions datatypes.JSONSlice[do.ModeratorPermission] `json:"requested_permissions"` // 请求的权限列表（外部传入）
}

// ── HTTP 请求结构体 ───────────────────────────────────────────────────────

// AddModeratorRequest 添加版主请求（管理员直接添加）
type AddModeratorRequest struct {
	UserID      uint                     `json:"user_id"      binding:"required" example:"1"`
	BoardID     uint                     `json:"board_id"     binding:"required" example:"1"`
	Permissions []do.ModeratorPermission `json:"permissions"` // 授予的权限列表
}

// Validate 校验
func (r *AddModeratorRequest) Validate() error {
	for _, perm := range r.Permissions {
		if !perm.IsValid() {
			return fmt.Errorf("无效的权限: %s", perm)
		}
	}
	return nil
}

// 更新权限请求
type UpdateModeratorPermissionsRequest struct {
	UserID      uint                                        `json:"user_id"`     // 申请人ID
	BoardID     uint                                        `json:"board_id"`    // 版块ID
	Permissions datatypes.JSONSlice[do.ModeratorPermission] `json:"permissions"` // 授予的权限列表
}

// BanUserRequest 封禁用户请求
type BanUserRequest struct {
	UserID    uint   `json:"user_id"                binding:"required" example:"1"`
	Reason    string `json:"reason"                 binding:"required" example:"发布违规内容"`
	ExpiresAt string `json:"expires_at" example:"2024-12-31T23:59:59Z"` // RFC3339
}

// ── 版主申请相关 ─────────────────────────────────────────────────────────

// ApplyModeratorInput 用户申请版主的参数（service 层使用）
// type ApplyModeratorInput struct {
// 	UserID               uint                     `json:"-"` // 从上下文注入
// 	BoardID              uint                     `json:"-"` // 从 URL 注入
// 	Reason               string                   `json:"reason" binding:"required,max=500"`
// 	RequestedPermissions []do.ModeratorPermission `json:"requested_permissions"` // 申请的权限列表
// }

// Validate 校验输入合法性
func (i *ApplyModeratorRequest) Validate() error {
	if i.Reason == "" {
		return fmt.Errorf("申请理由不能为空")
	}
	if len(i.RequestedPermissions) > 20 { // 防止恶意超大数组
		return fmt.Errorf("请求权限数量过多")
	}
	seen := make(map[do.ModeratorPermission]bool)
	for _, perm := range i.RequestedPermissions {
		if !perm.IsValid() {
			return fmt.Errorf("无效的权限: %s", perm)
		}
		if seen[perm] {
			return fmt.Errorf("权限不能重复: %s", perm)
		}
		seen[perm] = true
	}
	return nil
}
