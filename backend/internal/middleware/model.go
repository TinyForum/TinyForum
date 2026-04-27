package middleware

import (
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/repository/board"

	"github.com/gin-gonic/gin"
)

type MiddlewareSet struct {
	AuthMW               func() gin.HandlerFunc                                // 验证用户身份
	OptionalAuthMW       func() gin.HandlerFunc                                // 可选验证用户身份
	AdminRequiredMW      func() gin.HandlerFunc                                // 验证用户是否为管理员
	RateLimitMW          func(action ratelimit.Action) gin.HandlerFunc         // 限流
	ContentCheckMW       func(fields []string) gin.HandlerFunc                 // 内容检查
	ModeratorRequiredMW  func(boardRepo board.BoardRepository) gin.HandlerFunc // 验证用户是否为版主
	CanManageModeratorMW func(boardRepo board.BoardRepository) gin.HandlerFunc // 验证用户是否可以管理版主
	CanBanUserMW         func(boardRepo board.BoardRepository) gin.HandlerFunc // 验证用户是否可以封禁用户
	CanDeletePostMW      func(boardRepo board.BoardRepository) gin.HandlerFunc // 验证用户是否可以删除帖子
	CanPinPostMW         func(boardRepo board.BoardRepository) gin.HandlerFunc // 验证用户是否可以置顶帖子
}
