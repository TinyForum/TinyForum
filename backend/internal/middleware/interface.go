// internal/middleware/model.go
//
// 变更说明（相对原版）：
//   1. middlewareSet 新增 enforcer *casbin.Enforcer 字段
//   2. MiddlewareSet 接口新增 CasbinAuth() 方法
//   3. NewMiddlewareSet 新增 enforcer 参数
//   4. 修复 permission.go 中的类型断言隐患：
//      Auth 中间件向 context 注入的 user_role 是 string（来自 JWT claims.Role），
//      而原 RequirePermission 用 do.UserRole 断言，必然失败。
//      解决方案：在 middlewareSet 内部统一转换，外部调用者无感知。

package middleware

import (
	"tiny-forum/config"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/repository/board"
	"tiny-forum/internal/repository/token"
	"tiny-forum/internal/service/check"
	riskservice "tiny-forum/internal/service/risk"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MiddlewareSet 中间件集合接口
type MiddlewareSet interface {
	Auth() gin.HandlerFunc
	OptionalAuth() gin.HandlerFunc
	AdminRequired() gin.HandlerFunc
	// CasbinAuth 路由级 RBAC 鉴权（替代硬编码角色判断）
	// 需搭配 Auth() 或 OptionalAuth() 使用，确保 user_role 已注入 context
	CasbinAuth() gin.HandlerFunc
	RateLimit(action ratelimit.Action) gin.HandlerFunc
	ContentCheck(fields []string) gin.HandlerFunc
	ModeratorRequired(boardRepo board.BoardRepository) gin.HandlerFunc
	CanManageModerator(boardRepo board.BoardRepository) gin.HandlerFunc
	CanBanUser(boardRepo board.BoardRepository) gin.HandlerFunc
	CanDeletePost(boardRepo board.BoardRepository) gin.HandlerFunc
	CanPinPost(boardRepo board.BoardRepository) gin.HandlerFunc
}

// middlewareSet 私有实现
type middlewareSet struct {
	jwtMgr          *jwtpkg.JWTManager
	db              *gorm.DB
	riskSvc         riskservice.RiskService
	contentCheckSvc check.ContentCheckService
	tokenRepo       token.TokenRepository
	rateLimitCfg    *config.RateLimitConfig
	enforcer        *casbin.Enforcer // Casbin enforcer，由 wire 层注入
}

// NewMiddlewareSet 创建中间件集合实例。
//
// enforcer 参数：由 casbinx.NewEnforcer 创建后从 wire 层传入。
// 若暂不启用 Casbin，可传 nil，CasbinAuth() 调用时会直接放行（降级行为）。
func NewMiddlewareSet(
	jwtMgr *jwtpkg.JWTManager,
	db *gorm.DB,
	riskSvc riskservice.RiskService,
	contentCheckSvc check.ContentCheckService,
	tokenRepo token.TokenRepository,
	rateLimitCfg *config.RateLimitConfig,
	enforcer *casbin.Enforcer, // 新增
) MiddlewareSet {
	return &middlewareSet{
		jwtMgr:          jwtMgr,
		db:              db,
		riskSvc:         riskSvc,
		contentCheckSvc: contentCheckSvc,
		tokenRepo:       tokenRepo,
		rateLimitCfg:    rateLimitCfg,
		enforcer:        enforcer,
	}
}

// ── 接口实现 ──────────────────────────────────────────────────────────────────

func (m *middlewareSet) Auth() gin.HandlerFunc {
	return Auth(m.jwtMgr, m.tokenRepo)
}

func (m *middlewareSet) OptionalAuth() gin.HandlerFunc {
	return OptionalAuth(m.jwtMgr)
}

func (m *middlewareSet) AdminRequired() gin.HandlerFunc {
	return AdminRequired()
}

// CasbinAuth 返回 Casbin 路由级鉴权中间件。
// 若 enforcer 未初始化（nil），则跳过鉴权直接放行，便于测试环境禁用。
func (m *middlewareSet) CasbinAuth() gin.HandlerFunc {
	if m.enforcer == nil {
		return func(c *gin.Context) { c.Next() }
	}
	return casbinAuth(m.enforcer)
}

func (m *middlewareSet) RateLimit(action ratelimit.Action) gin.HandlerFunc {
	rateLimitMW := NewRateLimitMiddleware(m.db, m.riskSvc, m.rateLimitCfg)
	return rateLimitMW.Middleware(action)
}

func (m *middlewareSet) ContentCheck(fields []string) gin.HandlerFunc {
	return ContentCheckMiddleware(m.contentCheckSvc, fields)
}

func (m *middlewareSet) ModeratorRequired(boardRepo board.BoardRepository) gin.HandlerFunc {
	return ModeratorRequired(m.jwtMgr, boardRepo)
}

func (m *middlewareSet) CanManageModerator(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanManageModerator(m.jwtMgr, boardRepo)
}

func (m *middlewareSet) CanBanUser(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanBanUser(m.jwtMgr, boardRepo)
}

func (m *middlewareSet) CanDeletePost(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanDeletePost(m.jwtMgr, boardRepo)
}

func (m *middlewareSet) CanPinPost(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanPinPost(m.jwtMgr, boardRepo)
}
