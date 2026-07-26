package topic

import (
	"testing"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// 创建内存数据库，并自动迁移所有需要的表（Topic 和 User）
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 启用外键约束
	db.Exec("PRAGMA foreign_keys = ON")

	// Only migrate tables directly needed; avoid cascading from TopicCreation → Creation → Board
	err = db.AutoMigrate(&do.Topic{}, &do.User{}, &do.TopicFollow{})
	require.NoError(t, err)

	// Create tables manually to avoid GORM cascade pulling in boards/articles/questions/posts
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
	db.Exec(`CREATE TABLE IF NOT EXISTS topic_creations (
		id INTEGER PRIMARY KEY,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		topic_id INTEGER NOT NULL,
		creation_id INTEGER NOT NULL,
		series_id INTEGER,
		viewpoint TEXT,
		is_public INTEGER DEFAULT 1,
		post_count INTEGER DEFAULT 0,
		follower_count INTEGER DEFAULT 0,
		sort_order INTEGER DEFAULT 0,
		creator_id INTEGER NOT NULL
	)`)

	return db
}

// 创建测试用的用户（方便关联）
func createTestUser(db *gorm.DB, id uint, username string) *do.User {
	user := &do.User{
		BaseModel: common.BaseModel{ID: id},
		Username:  username,
		Email:     username + "@test.com",
		Password:  "hashed",
	}
	db.Create(user)
	return user
}

// 创建测试话题
func createTestTopic(db *gorm.DB, id uint, title, slug string, creatorID uint) *do.Topic {
	topic := &do.Topic{
		BaseModel: common.BaseModel{ID: id},
		Title:     title,
		Slug:      slug,
		CreatorID: creatorID,
	}
	db.Create(topic)
	return topic
}

// 创建测试作品（Post / Creation）
func createTestCreation(db *gorm.DB, id uint, title string, authorID uint) *do.Creation {
	creation := &do.Creation{
		BaseModel: common.BaseModel{ID: id},
		Title:     title,
		AuthorID:  authorID,
		ImageUrls: []string{"https://example.com/image.jpg", "https://example.com/image2.jpg", "https://example.com/image3.jpg"},
	}
	db.Create(creation)
	return creation
}
