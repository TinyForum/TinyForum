package user

import (
	"context"
	"testing"
	"time"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupGlobalStatsDB 在通用统计测试库基础上补齐点赞/关注/收藏/违规表。
func setupGlobalStatsDB(t *testing.T) *gorm.DB {
	db := setupStatsTestDB(t)

	// GetGlobalStatsCount 使用 errgroup 并发查询，sqlite :memory: 每个连接是独立库，
	// 必须限制为单连接，否则其它 goroutine 会看不到已建的表。
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id INTEGER NOT NULL,
		target_type TEXT NOT NULL,
		target_id INTEGER NOT NULL
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS follows (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		follower_id INTEGER NOT NULL,
		following_id INTEGER NOT NULL
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS favorites (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id INTEGER NOT NULL,
		target_id INTEGER NOT NULL,
		target_type TEXT NOT NULL,
		group_id INTEGER NOT NULL
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS violations (
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
		status TEXT NOT NULL DEFAULT 'pending'
	)`).Error)

	return db
}

// TestGetGlobalStatsCount_NormalizedSchema 验证统计查询基于规范化后的 schema
// （creations.author_id / replies.author_id / likes.target_type+target_id）能正确计数。
func TestGetGlobalStatsCount_NormalizedSchema(t *testing.T) {
	db := setupGlobalStatsDB(t)
	repo := &userRepository{db: db}
	now := time.Now()

	author := &do.User{BaseModel: common.BaseModel{CreatedAt: now}, Username: "author", Email: "author@test.com", Score: 42}
	fan := &do.User{BaseModel: common.BaseModel{CreatedAt: now}, Username: "fan", Email: "fan@test.com"}
	require.NoError(t, db.Create(author).Error)
	require.NoError(t, db.Create(fan).Error)

	// 作者的一个帖子（article + creation）
	db.Exec(`INSERT INTO creations (id, created_at, author_id, title, content, slug, type) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		1, now, author.ID, "Post", "Content", "post-1", "article")
	db.Exec(`INSERT INTO articles (id, created_at, creation_id) VALUES (?, ?, ?)`, 1, now, 1)

	// 作者的一条评论（reply + comment）
	db.Exec(`INSERT INTO replies (id, created_at, author_id, content, target_type, target_id) VALUES (?, ?, ?, ?, ?, ?)`,
		1, now, author.ID, "Reply", "creation", 1)
	db.Exec(`INSERT INTO comments (id, created_at, works_id, reply_id) VALUES (?, ?, ?, ?)`, 1, now, 1, 1)

	// fan 给作者的帖子和评论各点了一个赞
	db.Exec(`INSERT INTO likes (created_at, user_id, target_type, target_id) VALUES (?, ?, ?, ?)`, now, fan.ID, "creation", 1)
	db.Exec(`INSERT INTO likes (created_at, user_id, target_type, target_id) VALUES (?, ?, ?, ?)`, now, fan.ID, "comment", 1)

	// fan 关注 author
	db.Exec(`INSERT INTO follows (created_at, follower_id, following_id) VALUES (?, ?, ?)`, now, fan.ID, author.ID)

	// author 收藏了一条内容
	db.Exec(`INSERT INTO favorites (created_at, user_id, target_id, target_type, group_id) VALUES (?, ?, ?, ?, ?)`, now, author.ID, 1, "creation", 0)

	// author 有一条违规
	db.Exec(`INSERT INTO violations (created_at, user_id, violation_type, reason) VALUES (?, ?, ?, ?)`, now, author.ID, "abuse", "test")

	stats, err := repo.GetGlobalStatsCount(context.Background(), author.ID)
	require.NoError(t, err)

	assert.Equal(t, int64(1), stats.TotalPost, "author 应有 1 篇帖子")
	assert.Equal(t, int64(1), stats.TotalComment, "author 应有 1 条评论")
	assert.Equal(t, int64(2), stats.TotalLike, "author 的帖子+评论共收到 2 个赞")
	assert.Equal(t, int64(1), stats.TotalFavorite, "author 有 1 个收藏")
	assert.Equal(t, int64(1), stats.TotalViolation, "author 有 1 条违规")
	assert.Equal(t, int64(1), stats.TotalFollower, "author 有 1 个粉丝")
	assert.Equal(t, int64(0), stats.TotalFollowing, "author 未关注任何人")
	assert.Equal(t, 42, stats.TotalScore, "author 积分为 42")
}
