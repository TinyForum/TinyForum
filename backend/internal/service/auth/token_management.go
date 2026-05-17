package auth

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// 验证 token
func (s *authService) ValidateResetToken(ctx context.Context, token string) (bool, error) {
	isVaildToken, err := s.authRepo.ValidateResetToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return isVaildToken, nil
}

func (s *authService) GetUserEmailByResetToken(ctx context.Context, token string) (string, error) {
	return s.authRepo.GetUserEmailByResetToken(ctx, token)
}

// RevokeToken  将 Token 加入黑名单，注销后无法再使用
// 写入 tokenRepo 的黑名单，中间件验证时检查黑名单
func (s *authService) RevokeToken(ctx context.Context, rawToken string) error {
	// 解析 token 获取 JTI 和过期时间，精确存储黑名单
	claims, err := s.jwtMgr.Parse(rawToken)
	if err != nil {
		// token 已无效，无需加入黑名单
		return nil
	}
	// 将 JTI 存入黑名单，TTL 与 token 剩余有效期一致，节省存储
	return s.tokenRepo.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time)
}
