package wire

import (
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/repository/board"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MiddlewareSet 持有中间件需要的依赖，方便统一创建
type MiddlewareSet struct {
	AuthMW               func() gin.HandlerFunc // 实际需要传入 jwtMgr，所以我们保留函数工厂
	OptionalAuthMW       func() gin.HandlerFunc
	AdminRequiredMW      func() gin.HandlerFunc
	RateLimitMW          func(action ratelimit.Action) gin.HandlerFunc
	ContentCheckMW       func(fields []string) gin.HandlerFunc
	ModeratorRequiredMW  func(boardRepo board.BoardRepository) gin.HandlerFunc
	CanManageModeratorMW func(boardRepo board.BoardRepository) gin.HandlerFunc
	CanBanUserMW         func(boardRepo board.BoardRepository) gin.HandlerFunc
	CanDeletePostMW      func(boardRepo board.BoardRepository) gin.HandlerFunc
	CanPinPostMW         func(boardRepo board.BoardRepository) gin.HandlerFunc
}

// NewMiddlewareSet 创建中间件工厂（需要依赖注入）
// 注意：部分中间件需要运行时传入 boardRepo 或 action，因此返回的是构造函数。
func NewMiddlewareSet(jwtMgr *jwtpkg.JWTManager, db *gorm.DB, services *Services) *MiddlewareSet {
	return &MiddlewareSet{
		AuthMW:          func() gin.HandlerFunc { return middleware.Auth(jwtMgr) },
		OptionalAuthMW:  func() gin.HandlerFunc { return middleware.OptionalAuth(jwtMgr) },
		AdminRequiredMW: func() gin.HandlerFunc { return middleware.AdminRequired() },
		RateLimitMW: func(action ratelimit.Action) gin.HandlerFunc {
			return middleware.RateLimitMiddleware(db, services.Risk, action)
		},
		ContentCheckMW: func(fields []string) gin.HandlerFunc {
			return middleware.ContentCheckMiddleware(services.ContentCheck, fields)
		},
		ModeratorRequiredMW: func(boardRepo board.BoardRepository) gin.HandlerFunc {
			return middleware.ModeratorRequired(jwtMgr, boardRepo)
		},
		CanManageModeratorMW: func(boardRepo board.BoardRepository) gin.HandlerFunc {
			return middleware.CanManageModerator(jwtMgr, boardRepo)
		},
		CanBanUserMW: func(boardRepo board.BoardRepository) gin.HandlerFunc {
			return middleware.CanBanUser(jwtMgr, boardRepo)
		},
		CanDeletePostMW: func(boardRepo board.BoardRepository) gin.HandlerFunc {
			return middleware.CanDeletePost(jwtMgr, boardRepo)
		},
		CanPinPostMW: func(boardRepo board.BoardRepository) gin.HandlerFunc {
			return middleware.CanPinPost(jwtMgr, boardRepo)
		},
	}
}
