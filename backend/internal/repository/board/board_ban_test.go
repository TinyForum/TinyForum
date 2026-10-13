package board

import (
	"testing"

	"tiny-forum/internal/model/do"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupBoardBanDB 通过 GORM AutoMigrate 建表，验证模型的唯一索引声明为"部分索引"（deleted_at IS NULL），
// 与 migrations/20261013_board_ban_partial_unique.up.sql 保持一致的语义。
func setupBoardBanDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&do.User{}, &do.BoardBan{}))
	return db
}

// TestBoardBan_ReBanAfterUnban 验证"封禁 → 解封(软删) → 再次封禁"在同一 (user_id, board_id) 上可成功。
// 若唯一索引未带 deleted_at 谓词，第二次 Create 会触发唯一约束冲突（回归该缺陷）。
func TestBoardBan_ReBanAfterUnban(t *testing.T) {
	db := setupBoardBanDB(t)
	repo := &boardRepository{db: db}

	ban := &do.BoardBan{UserID: 1, BoardID: 1, BannedBy: 2, Reason: "spam"}
	require.NoError(t, repo.BanUser(ban), "首次封禁应成功")

	require.NoError(t, repo.UnbanUser(1, 1), "解封应成功")

	// 关键回归点：同一用户在同一板块再次封禁不得报唯一约束错误
	require.NoError(t, repo.BanUser(&do.BoardBan{UserID: 1, BoardID: 1, BannedBy: 2, Reason: "again"}), "解封后重新封禁不应触发唯一索引冲突")

	banned, err := repo.IsBanned(1, 1)
	require.NoError(t, err)
	assert.True(t, banned, "重新封禁后用户应处于封禁状态")
}

// TestBoardBan_GetBanNotFound 验证未命中时不返回非法指针（旧实现返回 &ban + ErrRecordNotFound）。
func TestBoardBan_GetBanNotFound(t *testing.T) {
	db := setupBoardBanDB(t)
	repo := &boardRepository{db: db}

	ban, err := repo.GetBan(99, 99)
	require.NoError(t, err)
	assert.Nil(t, ban, "无记录时应返回 nil,nil")
}