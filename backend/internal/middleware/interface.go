// internal/middleware/model.go
package middleware

import (
	"tiny-forum/internal/infra/config"
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
	Auth() gin.HandlerFunc                                              // 需要登录
	OptionalAuth() gin.HandlerFunc                                      // 可选登录
	AdminRequired() gin.HandlerFunc                                     // 需要管理员权限
	SystemMaintainerRequired() gin.HandlerFunc                          // 需要系统维护者权限
	CasbinAuth() gin.HandlerFunc                                        // 路由鉴权
	RateLimit(action ratelimit.Action) gin.HandlerFunc                  // 限流
	ContentCheck(fields []string) gin.HandlerFunc                       // 内容合规
	ModeratorRequired(boardRepo board.BoardRepository) gin.HandlerFunc  // 需要版主权限
	CanManageModerator(boardRepo board.BoardRepository) gin.HandlerFunc // 可管理版主
	CanBanUser(boardRepo board.BoardRepository) gin.HandlerFunc         // 可封禁用户
	CanDeletePost(boardRepo board.BoardRepository) gin.HandlerFunc      // 可删除帖子
	CanPinPost(boardRepo board.BoardRepository) gin.HandlerFunc         // 可置顶帖子
	UpdateConfig(cfg *config.Config)                                    // 更新配置
}

// middlewareSet 私有实现
type middlewareSet struct {
	jwtMgr          *jwtpkg.JWTManager
	db              *gorm.DB
	riskSvc         riskservice.RiskService
	contentCheckSvc check.ContentCheckService
	tokenRepo       token.TokenRepository
	rateLimitCfg    *config.RateLimitConfig
	enforcer        *casbin.Enforcer
	// 保存动态配置引用
	dynCfg *config.DynamicConfig
	// 缓存中间件实例以便更新
	cachedRateMW *RateLimitMiddleware
}

// NewMiddlewareSet 创建中间件集合实例
func NewMiddlewareSet(
	jwtMgr *jwtpkg.JWTManager,
	db *gorm.DB,
	riskSvc riskservice.RiskService,
	contentCheckSvc check.ContentCheckService,
	tokenRepo token.TokenRepository,
	rateLimitCfg *config.RateLimitConfig,
	enforcer *casbin.Enforcer,
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

// NewMiddlewareSetWithDynamic 使用动态配置创建中间件集合
func NewMiddlewareSetWithDynamic(
	dynCfg *config.DynamicConfig,
	jwtMgr *jwtpkg.JWTManager,
	db *gorm.DB,
	riskSvc riskservice.RiskService,
	contentCheckSvc check.ContentCheckService,
	tokenRepo token.TokenRepository,
	enforcer *casbin.Enforcer,
) MiddlewareSet {
	cfg := dynCfg.Get()
	ms := &middlewareSet{
		jwtMgr:          jwtMgr,
		db:              db,
		riskSvc:         riskSvc,
		contentCheckSvc: contentCheckSvc,
		tokenRepo:       tokenRepo,
		rateLimitCfg:    &cfg.RiskControl.RateLimit,
		enforcer:        enforcer,
		dynCfg:          dynCfg,
	}

	// 注册配置变更回调
	dynCfg.OnChange(func(fileName string, oldConfig, newConfig *config.Config) {
		if fileName == "risk_control.yml" || fileName == "basic.yml" || fileName == "manual_reload" {
			ms.UpdateConfig(newConfig)
		}
	})

	return ms
}

// UpdateConfig 实现 MiddlewareSet 接口的 UpdateConfig 方法
func (m *middlewareSet) UpdateConfig(cfg *config.Config) {
	// 更新限流配置
	m.rateLimitCfg = &cfg.RiskControl.RateLimit

	// 如果缓存了限流中间件实例，更新其配置
	if m.cachedRateMW != nil {
		m.cachedRateMW.UpdateConfig(cfg.RiskControl.RateLimit)
	}

	// 更新其他中间件的配置
	// 注意：内容检查、版主权限等中间件可能也需要更新
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

// SystemMaintainerRequired 系统维护员权限中间件
// 只有系统维护员（SystemMaintainer）角色可以访问
func (m *middlewareSet) SystemMaintainerRequired() gin.HandlerFunc {
	return SystemMaintainerRequired()
}

func (m *middlewareSet) CasbinAuth() gin.HandlerFunc {
	if m.enforcer == nil {
		return func(c *gin.Context) { c.Next() }
	}
	return casbinAuth(m.enforcer)
}

func (m *middlewareSet) RateLimit(action ratelimit.Action) gin.HandlerFunc {
	// 创建限流中间件并缓存
	rateLimitMW := NewRateLimitMiddleware(m.db, m.riskSvc, m.rateLimitCfg)
	m.cachedRateMW = rateLimitMW
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
