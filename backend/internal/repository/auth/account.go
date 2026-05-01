package auth

import (
	"context"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

func (r *authRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&po.User{}, id).Error
}

// Restore 恢复软删除的账户（将 deleted_at 置空）
func (r *authRepository) Restore(ctx context.Context, id uint) error {
	return r.db.Unscoped().Model(&po.User{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// HardDelete 硬删除（物理删除）
func (r *authRepository) HardDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&po.User{}, id).Error
}

// HardDeleteWithTx 事务中硬删除
func (r *authRepository) HardDeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Unscoped().Delete(&po.User{}, id).Error
}

// GetDeletedUser 获取已软删除的用户
func (r *authRepository) GetDeletedUser(ctx context.Context, id uint) (*po.User, error) {
	var user po.User
	err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&user).Error
	return &user, err
}

// GetUserWithDeleted 获取用户信息（包括已软删除的）
func (r *authRepository) GetUserWithDeleted(ctx context.Context, id uint) (*po.User, error) {
	var user po.User
	// Unscoped() 可以查询到已软删除的记录
	err := r.db.WithContext(ctx).Unscoped().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
