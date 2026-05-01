// internal/auth/repository/interface.go
package auth

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*do.User, error)              // 通过邮箱查找用户
	Update(ctx context.Context, user *do.User) error                              // 更新用户信息
	Save(ctx context.Context, user *do.User) error                                // 保存用户信息
	SoftDelete(ctx context.Context, id uint) error                                // 软删除用户
	Restore(ctx context.Context, id uint) error                                   // 恢复软删除用户
	HardDelete(ctx context.Context, id uint) error                                // 硬删除用户
	HardDeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error             // 在事务中硬删除用户
	GetDeletedUser(ctx context.Context, id uint) (*do.User, error)                // 获取软删除用户
	GetUserWithDeleted(ctx context.Context, id uint) (*do.User, error)            // 获取包含软删除信息的用户
	GetUserByResetToken(ctx context.Context, token string) (*do.User, error)      // 通过重置密码令牌查找用户
	GetUserEmailByResetToken(ctx context.Context, token string) (string, error)   // 通过重置密码令牌查找用户
	MarkTokenAsUsed(ctx context.Context, tokenID uint) error                      // 标记令牌为已使用
	FindByResetToken(ctx context.Context, token string) (*do.RefreshToken, error) // 通过重置密码令牌查找用户
	ValidateResetToken(ctx context.Context, token string) (bool, error)           // 验证重置密码令牌
	DeleteResetToken(ctx context.Context, token string) error                     // 删除重置密码令牌

}
type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
