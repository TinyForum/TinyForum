package favorite

import (
	"context"
	"testing"

	"tiny-forum/internal/model/do"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupFavoriteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&do.Favorite{}))
	// creations / boards 仅用于联表查询断言，手工建最小表结构
	require.NoError(t, db.Exec(`CREATE TABLE creations (
		id INTEGER PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime,
		title text, summary text, like_count integer, board_id integer)`).Error)
	require.NoError(t, db.Exec(`CREATE TABLE boards (
		id INTEGER PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime, name text)`).Error)
	return db
}

// TestFavorite_ReFavoriteAfterUnfavorite 验证"收藏 → 取消(软删) → 再次收藏"同一目标可行。
// 若唯一索引未带 deleted_at 谓词，第二次 Create 会触发唯一约束冲突（回归该缺陷）。
func TestFavorite_ReFavoriteAfterUnfavorite(t *testing.T) {
	db := setupFavoriteDB(t)
	repo := NewFavoriteRepository(db)
	ctx := context.Background()

	fav := &do.Favorite{UserID: 1, TargetID: 10, TargetType: FavoriteTargetTypeCreation, GroupID: 0}
	require.NoError(t, repo.Create(ctx, fav), "首次收藏应成功")

	deleted, err := repo.Delete(ctx, 1, 10, FavoriteTargetTypeCreation)
	require.NoError(t, err)
	require.EqualValues(t, 1, deleted, "取消收藏应软删 1 行")

	require.NoError(t, repo.Create(ctx, &do.Favorite{UserID: 1, TargetID: 10, TargetType: FavoriteTargetTypeCreation, GroupID: 0}),
		"取消后重新收藏不应触发唯一约束冲突")
}

// TestFavorite_Exists 验证收藏存在性判断区分软删记录。
func TestFavorite_Exists(t *testing.T) {
	db := setupFavoriteDB(t)
	repo := NewFavoriteRepository(db)
	ctx := context.Background()

	exists, err := repo.Exists(ctx, 1, 10, FavoriteTargetTypeCreation)
	require.NoError(t, err)
	assert.False(t, exists)

	require.NoError(t, repo.Create(ctx, &do.Favorite{UserID: 1, TargetID: 10, TargetType: FavoriteTargetTypeCreation}))
	exists, err = repo.Exists(ctx, 1, 10, FavoriteTargetTypeCreation)
	require.NoError(t, err)
	assert.True(t, exists)

	_, err = repo.Delete(ctx, 1, 10, FavoriteTargetTypeCreation)
	require.NoError(t, err)
	exists, err = repo.Exists(ctx, 1, 10, FavoriteTargetTypeCreation)
	require.NoError(t, err)
	assert.False(t, exists, "软删后应视为未收藏")
}

// TestFavorite_ListByUser 验证收藏列表联表返回作品信息并过滤软删。
func TestFavorite_ListByUser(t *testing.T) {
	db := setupFavoriteDB(t)
	repo := NewFavoriteRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Exec(`INSERT INTO boards (id, name) VALUES (1, '技术')`).Error)
	require.NoError(t, db.Exec(`INSERT INTO creations (id, title, summary, like_count, board_id) VALUES (10, '帖子A', '摘要A', 5, 1)`).Error)
	require.NoError(t, db.Exec(`INSERT INTO creations (id, title, summary, like_count, board_id) VALUES (11, '帖子B', '摘要B', 2, 1)`).Error)

	require.NoError(t, repo.Create(ctx, &do.Favorite{UserID: 1, TargetID: 10, TargetType: FavoriteTargetTypeCreation}))
	require.NoError(t, repo.Create(ctx, &do.Favorite{UserID: 1, TargetID: 11, TargetType: FavoriteTargetTypeCreation}))
	// 他人收藏不应出现在用户 1 的列表中
	require.NoError(t, repo.Create(ctx, &do.Favorite{UserID: 2, TargetID: 11, TargetType: FavoriteTargetTypeCreation}))

	rows, total, err := repo.ListByUser(ctx, 1, FavoriteTargetTypeCreation, 0, 20)
	require.NoError(t, err)
	assert.EqualValues(t, 2, total)
	require.Len(t, rows, 2)
	for _, row := range rows {
		assert.Equal(t, "技术", row.BoardName)
		assert.NotEmpty(t, row.Title)
	}
}
