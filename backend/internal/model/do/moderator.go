package do

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"tiny-forum/internal/model/common"
)

// ── ModeratorPermission 版主权限标识（可动态扩展）────────────────────────────

type ModeratorPermission string

const (
	PerModDeletePost      ModeratorPermission = "delete_post"      // 删除帖子
	PerMoePinPost         ModeratorPermission = "pin_post"         // 置顶帖子
	PerModEditAnyPost     ModeratorPermission = "edit_any_post"    // 编辑任意帖子
	PerModManageModerator ModeratorPermission = "manage_moderator" // 管理版主
	PerModBanUser         ModeratorPermission = "ban_user"         // 封禁用户
	// 未来新增权限只需添加常量，无需改表结构
)

// 版主权限有效性检验
func (p ModeratorPermission) IsValid() bool {
	switch p {
	case PerModDeletePost, PerMoePinPost, PerModEditAnyPost, PerModManageModerator, PerModBanUser:
		return true
	}
	return false
}

// ParsePermission 严格解析，返回错误
func ParsePermission(s string) (ModeratorPermission, error) {
	p := ModeratorPermission(s)
	if p.IsValid() {
		return p, nil
	}
	return "", fmt.Errorf("invalid permission: %s", s)
}

// MarshalJSON 确保序列化为字符串
func (p ModeratorPermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}

// UnmarshalJSON 反序列化
func (p *ModeratorPermission) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*p = ModeratorPermission(s)
	return nil
}

// ── PermissionSet 权限集合（用于 JSON 存储）────────────────────────────────

// PermissionSet 代表一组权限，本质是权限切片，支持动态增删
type ModeratorPermissionSet []ModeratorPermission

// Value 实现 driver.Valuer 接口，用于写入数据库
func (ps ModeratorPermissionSet) Value() (driver.Value, error) {
	if ps == nil {
		return nil, nil
	}
	return json.Marshal(ps)
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取
func (ps *ModeratorPermissionSet) Scan(src interface{}) error {
	if src == nil {
		*ps = ModeratorPermissionSet{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for PermissionSet: %T", src)
	}
	return json.Unmarshal(data, ps)
}

// Contains 检查是否包含某个权限
func (ps ModeratorPermissionSet) Contains(perm ModeratorPermission) bool {
	for _, p := range ps {
		if p == perm {
			return true
		}
	}
	return false
}

// Add 添加权限（去重）
func (ps *ModeratorPermissionSet) Add(perm ModeratorPermission) {
	if !ps.Contains(perm) {
		*ps = append(*ps, perm)
	}
}

// Remove 移除权限
func (ps *ModeratorPermissionSet) Remove(perm ModeratorPermission) {
	for i, p := range *ps {
		if p == perm {
			*ps = append((*ps)[:i], (*ps)[i+1:]...)
			return
		}
	}
}

// ModeratorBoardInfo 用户管理的板块信息（含权限）
type ModeratorBoardWithPerms struct {
	Board
	Permissions ModeratorPermissionSet `gorm:"type:json" json:"permissions"`
}

// ── Moderator 板块版主记录 ─────────────────────────────────────────────────

// Moderator 板块版主记录，权限使用 PermissionSet 存储在 JSON 列中
type Moderator struct {
	common.BaseModel
	UserID      uint                   `gorm:"not null;uniqueIndex:idx_user_board" json:"user_id"`
	BoardID     uint                   `gorm:"not null;uniqueIndex:idx_user_board" json:"board_id"`
	Permissions ModeratorPermissionSet `gorm:"type:json" json:"permissions"`
	User        User                   `gorm:"foreignKey:UserID"  json:"user,omitempty"`
	Board       Board                  `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}

// HasPermission 检查版主是否拥有某权限
func (m *Moderator) HasPermission(perm ModeratorPermission) bool {
	return m.Permissions.Contains(perm)
}

// SetPermissions 替换整个权限集合
func (m *Moderator) SetPermissions(perms ModeratorPermissionSet) {
	m.Permissions = perms
}

// AddPermission 添加单个权限
func (m *Moderator) AddPermission(perm ModeratorPermission) {
	m.Permissions.Add(perm)
}

// RemovePermission 移除单个权限
func (m *Moderator) RemovePermission(perm ModeratorPermission) {
	m.Permissions.Remove(perm)
}

// 获取权限
func (m *Moderator) GetPermissions() (ModeratorPermissionSet, error) {

	return m.Permissions, nil
}

// ── 版主申请相关 ─────────────────────────────────────────────────────────

// ApplyModeratorInput 用户申请版主的参数（service 层使用）
type ApplyModeratorInput struct {
	UserID               uint                  `json:"-"` // 从上下文注入
	BoardID              uint                  `json:"-"` // 从 URL 注入
	Reason               string                `json:"reason" binding:"required,max=500"`
	RequestedPermissions []ModeratorPermission `json:"requested_permissions"` // 申请的权限列表
}

// Validate 校验输入合法性
func (i *ApplyModeratorInput) Validate() error {
	if i.Reason == "" {
		return fmt.Errorf("申请理由不能为空")
	}
	if len(i.RequestedPermissions) > 20 { // 防止恶意超大数组
		return fmt.Errorf("请求权限数量过多")
	}
	seen := make(map[ModeratorPermission]bool)
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
