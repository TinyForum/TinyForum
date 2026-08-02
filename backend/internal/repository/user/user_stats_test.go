package user

import (
	"context"
	"testing"
	"time"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupStatsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db.Exec("PRAGMA foreign_keys = ON")

	err = db.AutoMigrate(&do.User{})
	require.NoError(t, err)

	db.Exec(`CREATE TABLE IF NOT EXISTS creations (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		title TEXT NOT NULL,
		content TEXT NOT NULL DEFAULT '',
		summary TEXT,
		slug TEXT,
		cover_url TEXT,
		image_urls TEXT,
		type TEXT NOT NULL DEFAULT '',
		creation_status TEXT DEFAULT 'draft',
		moderation_status TEXT DEFAULT 'normal',
		author_id INTEGER NOT NULL,
		view_count INTEGER DEFAULT 0,
		like_count INTEGER DEFAULT 0,
		pin_top INTEGER DEFAULT 0,
		is_original INTEGER DEFAULT 1,
		board_id INTEGER,
		pin_in_board INTEGER DEFAULT 0
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		creation_id INTEGER NOT NULL UNIQUE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS replies (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		target_type TEXT NOT NULL DEFAULT '',
		target_id INTEGER NOT NULL,
		parent_id INTEGER,
		like_count INTEGER DEFAULT 0,
		status TEXT DEFAULT 'visible'
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		works_id INTEGER NOT NULL,
		reply_id INTEGER NOT NULL UNIQUE,
		is_pinned INTEGER DEFAULT 0,
		is_anonymous INTEGER DEFAULT 0,
		dislike_count INTEGER DEFAULT 0,
		report_count INTEGER DEFAULT 0,
		sort_weight INTEGER DEFAULT 0,
		ip_location TEXT
	)`)

	return db
}

func TestCountActiveByDateRange_SQLCorrectness(t *testing.T) {
	db := setupStatsTestDB(t)
	repo := &userRepository{db: db}

	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	user := &do.User{
		BaseModel: common.BaseModel{CreatedAt: now},
		Username:  "testuser",
		Email:     "test@test.com",
	}
	require.NoError(t, db.Create(user).Error)

	// 插入 creation 并关联 article
	db.Exec(`INSERT INTO creations (id, created_at, author_id, title, content, slug, type, board_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		1, now, user.ID, "Test", "Content", "test-slug", "post", 0)
	db.Exec(`INSERT INTO articles (id, created_at, creation_id) VALUES (?, ?, ?)`,
		1, now, 1)

	// 插入 reply 并关联 comment
	db.Exec(`INSERT INTO replies (id, created_at, author_id, content, target_type, target_id) VALUES (?, ?, ?, ?, ?, ?)`,
		1, now, user.ID, "Test Reply", "post", 1)
	db.Exec(`INSERT INTO comments (id, created_at, works_id, reply_id) VALUES (?, ?, ?, ?)`,
		1, now, 1, 1)

	// 这应该成功执行，不再报 "column p.author_id does not exist"
	count, err := repo.CountActiveByDateRange(context.Background(), startDate, endDate)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count, "应有1个活跃用户（同时发了文章和评论）")

	// 测试仅有评论的活跃用户
	user2 := &do.User{
		BaseModel: common.BaseModel{CreatedAt: now},
		Username:  "testuser2",
		Email:     "test2@test.com",
	}
	require.NoError(t, db.Create(user2).Error)

	db.Exec(`INSERT INTO replies (id, created_at, author_id, content, target_type, target_id) VALUES (?, ?, ?, ?, ?, ?)`,
		2, now, user2.ID, "Test Reply 2", "post", 2)
	db.Exec(`INSERT INTO comments (id, created_at, works_id, reply_id) VALUES (?, ?, ?, ?)`,
		2, now, 2, 2)

	count, err = repo.CountActiveByDateRange(context.Background(), startDate, endDate)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count, "应有2个活跃用户")
}

func TestGetActiveUsersByDateRange_SQLCorrectness(t *testing.T) {
	db := setupStatsTestDB(t)
	repo := &userRepository{db: db}

	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	user := &do.User{
		BaseModel: common.BaseModel{CreatedAt: now},
		Username:  "activeuser",
		Email:     "active@test.com",
	}
	require.NoError(t, db.Create(user).Error)

	db.Exec(`INSERT INTO creations (id, created_at, author_id, title, content, slug, type, board_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		1, now, user.ID, "Test", "Content", "active-slug", "post", 0)
	db.Exec(`INSERT INTO articles (id, created_at, creation_id) VALUES (?, ?, ?)`,
		1, now, 1)

	// 应成功返回活跃用户列表
	rows, err := repo.GetActiveUsersByDateRange(context.Background(), startDate, endDate, 10)
	assert.NoError(t, err)
	assert.NotNil(t, rows)
}
