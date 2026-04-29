// internal/middleware/middleware_set.go

package middleware

import (
	"tiny-forum/config"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/repository/board"
	"tiny-forum/internal/repository/token"
	"tiny-forum/internal/service/check"
	riskservice "tiny-forum/internal/service/risk"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MiddlewareSet 中间件集合接口
type MiddlewareSet interface {
	Auth() gin.HandlerFunc
	OptionalAuth() gin.HandlerFunc
	AdminRequired() gin.HandlerFunc
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
	// 依赖
	jwtMgr          *jwtpkg.JWTManager
	db              *gorm.DB
	riskSvc         riskservice.RiskService
	contentCheckSvc check.ContentCheckService
	tokenRepo       token.TokenRepository
	rateLimitCfg    *config.RateLimitConfig
}

// NewMiddlewareSet 创建中间件集合实例
func NewMiddlewareSet(
	jwtMgr *jwtpkg.JWTManager,
	db *gorm.DB,
	riskSvc riskservice.RiskService,
	contentCheckSvc check.ContentCheckService,
	tokenRepo token.TokenRepository,
	rateLimitCfg *config.RateLimitConfig,
) MiddlewareSet {
	return &middlewareSet{
		jwtMgr:          jwtMgr,
		db:              db,
		riskSvc:         riskSvc,
		contentCheckSvc: contentCheckSvc,
		tokenRepo:       tokenRepo,
		rateLimitCfg:    rateLimitCfg,
	}
}

// Auth 验证用户身份
func (m *middlewareSet) Auth() gin.HandlerFunc {
	return Auth(m.jwtMgr, m.tokenRepo)
}

// OptionalAuth 可选验证用户身份
func (m *middlewareSet) OptionalAuth() gin.HandlerFunc {
	return OptionalAuth(m.jwtMgr)
}

// AdminRequired 验证用户是否为管理员
func (m *middlewareSet) AdminRequired() gin.HandlerFunc {
	return AdminRequired()
}

// RateLimit 限流
func (m *middlewareSet) RateLimit(action ratelimit.Action) gin.HandlerFunc {
	rateLimitMW := NewRateLimitMiddleware(m.db, m.riskSvc, m.rateLimitCfg)
	return rateLimitMW.Middleware(action)
}

// ContentCheck 内容检查
func (m *middlewareSet) ContentCheck(fields []string) gin.HandlerFunc {
	return ContentCheckMiddleware(m.contentCheckSvc, fields)
}

// ModeratorRequired 验证用户是否为版主
func (m *middlewareSet) ModeratorRequired(boardRepo board.BoardRepository) gin.HandlerFunc {
	return ModeratorRequired(m.jwtMgr, boardRepo)
}

// CanManageModerator 验证用户是否可以管理版主
func (m *middlewareSet) CanManageModerator(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanManageModerator(m.jwtMgr, boardRepo)
}

// CanBanUser 验证用户是否可以封禁用户
func (m *middlewareSet) CanBanUser(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanBanUser(m.jwtMgr, boardRepo)
}

// CanDeletePost 验证用户是否可以删除帖子
func (m *middlewareSet) CanDeletePost(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanDeletePost(m.jwtMgr, boardRepo)
}

// CanPinPost 验证用户是否可以置顶帖子
func (m *middlewareSet) CanPinPost(boardRepo board.BoardRepository) gin.HandlerFunc {
	return CanPinPost(m.jwtMgr, boardRepo)
}