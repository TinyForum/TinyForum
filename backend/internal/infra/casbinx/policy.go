// internal/infra/casbinx/policy.go
//
// PolicyManager 封装运行时策略管理操作，供 admin handler 调用。
// 所有写操作都会立即 SavePolicy 持久化到数据库，无需重启生效。

package casbinx

import (
	"fmt"

	"github.com/casbin/casbin/v3"
)

// PolicyManager 运行时策略管理
type PolicyManager struct {
	enforcer *casbin.Enforcer
}

// NewPolicyManager 创建策略管理器
func NewPolicyManager(enforcer *casbin.Enforcer) *PolicyManager {
	return &PolicyManager{enforcer: enforcer}
}

// ── 路由策略 CRUD ─────────────────────────────────────────────────────────────

// AddPolicy 为角色添加路由权限。
//
//	role:   角色字符串，如 "moderator"
//	path:   路由路径，如 "/api/v1/boards/:id/ban"
//	method: HTTP 方法或正则，如 "POST" 或 "GET|POST"
func (m *PolicyManager) AddPolicy(role, path, method string) error {
	ok, err := m.enforcer.AddPolicy(role, path, method)
	if err != nil {
		return fmt.Errorf("AddPolicy: %w", err)
	}
	if !ok {
		return nil // 策略已存在，忽略
	}
	return m.enforcer.SavePolicy()
}

// RemovePolicy 删除角色的路由权限
func (m *PolicyManager) RemovePolicy(role, path, method string) error {
	ok, err := m.enforcer.RemovePolicy(role, path, method)
	if err != nil {
		return fmt.Errorf("RemovePolicy: %w", err)
	}
	if !ok {
		return nil // 策略不存在，忽略
	}
	return m.enforcer.SavePolicy()
}

// GetPoliciesForRole 查询某角色的所有策略
func (m *PolicyManager) GetPoliciesForRole(role string) ([][]string, error) {
	return m.enforcer.GetFilteredPolicy(0, role)
}

// ── 角色继承管理 ──────────────────────────────────────────────────────────────

// AddRoleInheritance 添加角色继承关系：child 继承 parent 的所有策略
func (m *PolicyManager) AddRoleInheritance(child, parent string) error {
	ok, err := m.enforcer.AddGroupingPolicy(child, parent)
	if err != nil {
		return fmt.Errorf("AddGroupingPolicy: %w", err)
	}
	if !ok {
		return nil
	}
	return m.enforcer.SavePolicy()
}

// RemoveRoleInheritance 删除角色继承关系
func (m *PolicyManager) RemoveRoleInheritance(child, parent string) error {
	ok, err := m.enforcer.RemoveGroupingPolicy(child, parent)
	if err != nil {
		return fmt.Errorf("RemoveGroupingPolicy: %w", err)
	}
	if !ok {
		return nil
	}
	return m.enforcer.SavePolicy()
}

// ── 用户级临时权限（特殊场景）────────────────────────────────────────────────
//
// 注意：TinyForum 的 JWT 中 sub 是角色字符串，不是用户 ID。
// 如需给单个用户赋予临时权限，应在 JWT 中携带特殊角色，
// 而不是直接在 Casbin 中绑定 user_id → policy。
//
// 以下方法仅供有此需求时参考，常规场景不需要使用。

// AssignRoleToUser 将用户（通过角色字符串标识）分配到某角色组
// subject 通常是 "user:<id>" 格式
func (m *PolicyManager) AssignRoleToUser(subject, role string) error {
	ok, err := m.enforcer.AddGroupingPolicy(subject, role)
	if err != nil {
		return fmt.Errorf("AssignRoleToUser: %w", err)
	}
	if !ok {
		return nil
	}
	return m.enforcer.SavePolicy()
}

// RemoveRoleFromUser 移除用户的角色分配
func (m *PolicyManager) RemoveRoleFromUser(subject, role string) error {
	ok, err := m.enforcer.RemoveGroupingPolicy(subject, role)
	if err != nil {
		return fmt.Errorf("RemoveRoleFromUser: %w", err)
	}
	if !ok {
		return nil
	}
	return m.enforcer.SavePolicy()
}

// ── 查询 ──────────────────────────────────────────────────────────────────────

// Enforce 手动检查权限（供 service 层在非 HTTP 场景使用）
func (m *PolicyManager) Enforce(role, path, method string) (bool, error) {
	return m.enforcer.Enforce(role, path, method)
}

// GetAllRoles 获取所有已定义的角色
func (m *PolicyManager) GetAllRoles() ([]string, error) {
	return m.enforcer.GetAllRoles()
}

// ReloadPolicy 从数据库重新加载策略（用于多实例部署时同步策略变更）
func (m *PolicyManager) ReloadPolicy() error {
	return m.enforcer.LoadPolicy()
}
