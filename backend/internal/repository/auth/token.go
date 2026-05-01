package auth

import (
	"context"
	"tiny-forum/internal/model/po"
)

func (r *authRepository) DeleteResetToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Where("jti LIKE ?", "reset_%").
		Delete(&po.RefreshToken{}).Error
}
