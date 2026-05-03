// internal/infra/casbinx/enforcer.go
//
// 职责：
//   1. 用 GORM adapter 把策略持久化到 casbin_rule 表（自动建表）
//   2. 首次启动时写入默认策略（幂等）
//   3. 暴露 NewEnforcer 供 wire 层调用
//
// 策略分两层：
//   - 路由级 RBAC：Casbin 管（本文件）
//   - 版主细粒度权限：保持 moderator.go 的 DB JSON 查询方式，Casbin 不介入

package casbinx

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewEnforcer 创建并初始化 Casbin enforcer。
//
//   - modelPath: rbac_model.conf 的文件路径（相对于工作目录）
//   - db:        已连接的 GORM 实例，adapter 会在其中自动创建 casbin_rule 表
func NewEnforcer(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	// gorm-adapter v3 默认使用传入 db 的数据库，表名 casbin_rule
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("casbin: create adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin: create enforcer: %w", err)
	}

	// 从数据库加载已有策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("casbin: load policy: %w", err)
	}

	// 写入默认策略（已存在的条目会被 casbin 自动去重，不会重复插入）
	if err := seedDefaultPolicies(enforcer); err != nil {
		return nil, fmt.Errorf("casbin: seed policy: %w", err)
	}

	return enforcer, nil
}

// ── 默认策略 ──────────────────────────────────────────────────────────────────
//
// 规则格式：p, <角色>, <路径通配符>, <HTTP方法正则>
//
// 路径使用 keyMatch2（* 匹配单段，** 匹配多段）
// 方法使用 regexMatch（支持 GET|POST 等 OR 写法，".*" 匹配所有方法）

// defaultPolicies 路由级权限策略表
// [sub, obj, act]
var defaultPolicies = [][]string{
	// ── guest（未登录） ───────────────────────────────────────────────────────
	// 只允许读取公开内容，不需要身份
	{"guest", "/api/v1/posts", "GET"},
	{"guest", "/api/v1/posts/*", "GET"},
	{"guest", "/api/v1/boards", "GET"},
	{"guest", "/api/v1/boards/*", "GET"},
	{"guest", "/api/v1/users/*", "GET"},
	{"guest", "/api/v1/topics", "GET"},
	{"guest", "/api/v1/topics/*", "GET"},
	{"guest", "/api/v1/tags", "GET"},
	{"guest", "/api/v1/tags/*", "GET"},
	{"guest", "/api/v1/announcements", "GET"},
	{"guest", "/api/v1/announcements/*", "GET"},
	{"guest", "/api/v1/stats/*", "GET"},
	{"guest", "/api/v1/health", "GET"},

	// ── user（已登录普通用户）────────────────────────────────────────────────
	// 继承 guest 的所有 GET 权限（通过 role_inheritance 实现）
	// 以下仅列写操作
	{"user", "/api/v1/posts", "POST"},
	{"user", "/api/v1/posts/:id", "PUT|DELETE"},
	{"user", "/api/v1/posts/:id/like", "POST|DELETE"},
	{"user", "/api/v1/comments", "POST"},
	{"user", "/api/v1/comments/:id", "PUT|DELETE"},
	{"user", "/api/v1/users/:id/follow", "POST|DELETE"},
	{"user", "/api/v1/notifications", "GET"},
	{"user", "/api/v1/notifications/*", "GET|PUT"},
	{"user", "/api/v1/timelines/*", "GET"},
	{"user", "/api/v1/topics/:id/follow", "POST|DELETE"},
	{"user", "/api/v1/boards/:id/apply", "POST"},
	{"user", "/api/v1/questions", "GET|POST"},
	{"user", "/api/v1/questions/*", "GET|POST|PUT|DELETE"},
	{"user", "/api/v1/answers", "GET|POST"},
	{"user", "/api/v1/answers/*", "GET|POST|PUT|DELETE"},
	{"user", "/api/v1/upload", "POST"},
	{"user", "/api/v1/auth/logout", "POST"},
	{"user", "/api/v1/auth/password", "PUT"},
	{"user", "/api/v1/auth/me", "GET|PUT"},

	// ── member（付费会员）──────────────────────────────────────────────────
	// 继承 user，额外权限在此添加
	// 目前与 user 路由相同，付费功能待扩展

	// ── reviewer（审核员）──────────────────────────────────────────────────
	// 继承 user，可访问审核相关路由
	{"reviewer", "/api/v1/admin/posts/pending", "GET"},
	{"reviewer", "/api/v1/admin/audit/tasks/:id/approve", "PUT"},
	{"reviewer", "/api/v1/admin/audit/tasks/:id/reject", "PUT"},

	// ── moderator（版主）────────────────────────────────────────────────────
	// 继承 user，版主级路由
	// 注：版主的细粒度板块权限（delete_post/pin_post 等）由 moderator.go 中间件负责
	//     Casbin 仅保证"版主角色可以访问这些路由"
	{"moderator", "/api/v1/boards/:id/ban", "POST|DELETE"},
	{"moderator", "/api/v1/boards/:id/posts/:post_id", "DELETE"},
	{"moderator", "/api/v1/boards/:id/posts/:post_id/pin", "PUT"},
	{"moderator", "/api/v1/boards/:id/moderators", "GET"},

	// ── admin（管理员）──────────────────────────────────────────────────────
	// 继承 moderator，可访问 /admin/* 下所有路由
	{"admin", "/api/v1/admin/*", ".*"},
	// admin 可以管理所有板块，不受板块 ID 限制
	{"admin", "/api/v1/boards", "GET|POST|PUT|DELETE"},
	{"admin", "/api/v1/boards/*", ".*"},
	// admin 可以管理所有用户，不受用户 ID 限制
	{"admin", "/api/v1/users", "GET|POST|PUT|DELETE"},
	{"admin", "/api/v1/users/*", ".*"},
	// admin 可以管理所有公告，不受公告 ID 限制
	{"admin", "/api/v1/announcements", "GET|POST|PUT|DELETE"},
	{"admin", "/api/v1/announcements/*", ".*"},

	// admin 可以查看所有统计数据
	{"admin", "/api/v1/stats", "GET"},
	{"admin", "/api/v1/stats/*", ".*"},

	// ── super_admin ──────────────────────────────────────────────────────────
	// 继承 admin，无额外路由限制（通配兜底）
	{"super_admin", "/api/v1/*", ".*"},
}

// roleInheritance 角色继承关系：[子角色, 父角色]
// 子角色自动拥有父角色的所有策略
var roleInheritance = [][]string{
	{"member", "user"},       // 会员继承普通用户
	{"reviewer", "user"},     // 审核员继承普通用户
	{"moderator", "user"},    // 版主继承普通用户
	{"admin", "moderator"},   // 管理员继承版主
	{"admin", "reviewer"},    // 管理员继承审核员
	{"super_admin", "admin"}, // 超管继承管理员
}

// seedDefaultPolicies 幂等地写入默认策略。
// casbin AddPolicy / AddGroupingPolicy 若条目已存在会返回 false 但不报错，可安全重复调用。
func seedDefaultPolicies(e *casbin.Enforcer) error {
	// 写入路由策略
	for _, p := range defaultPolicies {
		if _, err := e.AddPolicy(p[0], p[1], p[2]); err != nil {
			return fmt.Errorf("AddPolicy %v: %w", p, err)
		}
	}

	// 写入角色继承
	for _, g := range roleInheritance {
		if _, err := e.AddGroupingPolicy(g[0], g[1]); err != nil {
			return fmt.Errorf("AddGroupingPolicy %v: %w", g, err)
		}
	}

	// 持久化到数据库
	if err := e.SavePolicy(); err != nil {
		return fmt.Errorf("SavePolicy: %w", err)
	}

	return nil
}
