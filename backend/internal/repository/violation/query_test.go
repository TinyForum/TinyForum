package violation

import (
	"context"
	"testing"

	"tiny-forum/internal/model/request"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupViolationDB 构造最小 violations 表（与全局统计测试表结构保持一致）。
func setupViolationDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.Exec(`CREATE TABLE violations (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id INTEGER NOT NULL,
		violation_type TEXT NOT NULL,
		reason TEXT NOT NULL,
		content_snapshot TEXT,
		evidence_url TEXT,
		source TEXT NOT NULL DEFAULT 'system',
		status TEXT NOT NULL DEFAULT 'pending',
		operator_id INTEGER DEFAULT 0,
		punish_type TEXT DEFAULT '',
		punish_expire_at DATETIME,
		appeal_status TEXT NOT NULL DEFAULT 'none',
		appeal_reason TEXT,
		appeal_time DATETIME,
		appeal_result TEXT
	)`).Error)
	return db
}

// TestListAndCountUserViolation_WithFilters 验证：
// 1) CountUserViolationByUserID 返回真实总数（而非页内条数）；
// 2) status / keyword 过滤对列表与计数同时生效。
func TestListAndCountUserViolation_WithFilters(t *testing.T) {
	db := setupViolationDB(t)
	repo := &violationRepository{db: db}

	db.Exec(`INSERT INTO violations (user_id, violation_type, reason, status) VALUES (1, 'abuse', '辱骂其他用户', 'pending')`)
	db.Exec(`INSERT INTO violations (user_id, violation_type, reason, status) VALUES (1, 'spam', '刷屏', 'pending')`)
	db.Exec(`INSERT INTO violations (user_id, violation_type, reason, status) VALUES (1, 'ad', '推广广告', 'confirmed')`)
	db.Exec(`INSERT INTO violations (user_id, violation_type, reason, status) VALUES (2, 'spam', '他人违规', 'pending')`)

	// 无过滤：用户 1 共 3 条，用户 2 共 1 条
	reqAll := request.ListUserViolationRequest{}
	total, err := repo.CountUserViolationByUserID(context.Background(), reqAll, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total, "用户1 总违规数应为 3（非页内条数）")

	// 状态过滤：pending 且关键词命中 reason
	reqFilter := request.ListUserViolationRequest{Status: "pending", Keyword: "刷屏"}
	total, err = repo.CountUserViolationByUserID(context.Background(), reqFilter, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total, "状态=pending 且关键词=刷屏 应只命中 1 条")

	list, err := repo.ListUserViolationByUserID(context.Background(), reqFilter, 1)
	require.NoError(t, err)
	assert.Len(t, list, 1, "列表应与计数一致")
	assert.Equal(t, "刷屏", list[0].Reason)

	// 分页：page_size=2 时列表只返回 2 条，但计数仍为 3
	reqPage := request.ListUserViolationRequest{}
	reqPage.Page = 1
	reqPage.PageSize = 2
	total, err = repo.CountUserViolationByUserID(context.Background(), reqPage, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total, "分页不应影响计数")
	list, err = repo.ListUserViolationByUserID(context.Background(), reqPage, 1)
	require.NoError(t, err)
	assert.Len(t, list, 2, "page_size=2 应返回 2 条")
}