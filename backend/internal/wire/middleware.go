package wire

import (
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/repository/board"
	"tiny-forum/internal/repository/token"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewMiddlewareSet 创建中间件工厂（需要依赖注入）
func NewMiddlewareSet(
	jwtMgr *jwtpkg.JWTManager,
	db *gorm.DB,
	services *Services,
	tokenRepo token.TokenRepository) *middleware.MiddlewareSet {
	return &middleware.MiddlewareSet{
		AuthMW:          func() gin.HandlerFunc { return middleware.Auth(jwtMgr, tokenRepo) },
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
