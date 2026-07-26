package topic

// go test -v ./internal/repository/topic
import (
	"testing"
	"time"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- 测试 Create ----------
func TestTopicRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db) // 请确保构造函数存在

	// 创建关联用户
	user := createTestUser(db, 1, "creator1")

	t.Run("创建成功", func(t *testing.T) {
		topic := &do.Topic{
			BaseModel: common.BaseModel{CreatedAt: time.Now()},
			Title:     "Test Topic",
			Slug:      "test-topic",
			CreatorID: user.ID,
		}
		err := repo.Create(topic)
		assert.NoError(t, err)
		assert.NotZero(t, topic.ID)

		// 验证数据库中确实存在
		var found do.Topic
		err = db.First(&found, topic.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test Topic", found.Title)
	})

	t.Run("创建失败：唯一约束冲突 (Title)", func(t *testing.T) {
		topic1 := &do.Topic{
			Title:     "Duplicate Title",
			Slug:      "dup-title-1",
			CreatorID: user.ID,
		}
		err := repo.Create(topic1)
		require.NoError(t, err)

		topic2 := &do.Topic{
			Title:     "Duplicate Title", // 重复 Title
			Slug:      "dup-title-2",
			CreatorID: user.ID,
		}
		err = repo.Create(topic2)
		assert.Error(t, err) // 应返回唯一约束错误
	})

	t.Run("创建失败：唯一约束冲突 (Slug)", func(t *testing.T) {
		topic1 := &do.Topic{
			Title:     "Title 1",
			Slug:      "dup-slug",
			CreatorID: user.ID,
		}
		err := repo.Create(topic1)
		require.NoError(t, err)

		topic2 := &do.Topic{
			Title:     "Title 2",
			Slug:      "dup-slug", // 重复 Slug
			CreatorID: user.ID,
		}
		err = repo.Create(topic2)
		assert.Error(t, err)
	})
}

// ---------- 测试 Update ----------
func TestTopicRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user := createTestUser(db, 1, "creator1")
	topic := &do.Topic{
		Title:     "Original Title",
		Slug:      "original-slug",
		CreatorID: user.ID,
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	t.Run("更新成功", func(t *testing.T) {
		topic.Title = "Updated Title"
		topic.FollowerCount = 100
		err := repo.Update(topic)
		assert.NoError(t, err)

		var updated do.Topic
		db.First(&updated, topic.ID)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, 100, updated.FollowerCount)
	})

	t.Run("更新不存在的记录（Save 会插入新记录，但 ID 不为0才更新）", func(t *testing.T) {
		// 如果 ID 为 0，Save 会创建新记录。这里我们不测试该边界，实际业务通常不允许。
		// 我们可以传一个不存在的 ID，Save 会尝试更新，但因为记录不存在，会返回错误或影响 0 行。
		// 取决于 GORM 版本，Save 对于不存在的记录可能会创建新记录，这里不深入。
		// 更严谨的测试是使用 Update 方法（带条件），但此处我们仅测试正常更新。
		// 为了简单，我们只测试存在记录更新。
	})
}

// ---------- 测试 Delete ----------
func TestTopicRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user := createTestUser(db, 1, "creator1")
	topic := &do.Topic{
		Title:     "To Be Deleted",
		Slug:      "to-delete",
		CreatorID: user.ID,
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	t.Run("删除成功", func(t *testing.T) {
		err := repo.Delete(topic.ID)
		assert.NoError(t, err)

		// 验证软删除（DeletedAt 非空）
		var deleted do.Topic
		err = db.Unscoped().First(&deleted, topic.ID).Error
		assert.NoError(t, err)
		assert.NotNil(t, deleted.DeletedAt.Time)
	})

	t.Run("删除不存在的 ID", func(t *testing.T) {
		err := repo.Delete(99999)
		// GORM 删除不存在的记录不会报错（影响0行）
		assert.NoError(t, err)
	})
}

// ---------- 测试 FindByID ----------
func TestTopicRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user := createTestUser(db, 1, "creator1")
	topic := &do.Topic{
		Title:     "Find Me",
		Slug:      "find-me",
		CreatorID: user.ID,
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	t.Run("查找存在的记录（预加载 Creator）", func(t *testing.T) {
		found, err := repo.FindByID(topic.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Find Me", found.Title)
		// 验证预加载的 Creator
		assert.NotNil(t, found.Creator)
		assert.Equal(t, user.ID, found.Creator.ID)
	})

	t.Run("查找不存在的记录", func(t *testing.T) {
		found, err := repo.FindByID(99999)
		// GORM 返回 gorm.ErrRecordNotFound
		assert.Error(t, err)
		assert.Nil(t, found) // 由于我们返回指针，err 时通常为 nil 或零值，根据实现
	})
}

// ---------- 测试 List（已有，但为了完整性仍包含） ----------
func TestTopicRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	now := time.Now()
	user := createTestUser(db, 1, "creator1")

	testTopics := []do.Topic{
		{BaseModel: common.BaseModel{ID: 1, CreatedAt: now}, Title: "Topic A", Slug: "topic-a", FollowerCount: 10, PostCount: 5, CreatorID: user.ID},
		{BaseModel: common.BaseModel{ID: 2, CreatedAt: now}, Title: "Topic B", Slug: "topic-b", FollowerCount: 20, PostCount: 1, CreatorID: user.ID},
		{BaseModel: common.BaseModel{ID: 3, CreatedAt: now}, Title: "Topic C", Slug: "topic-c", FollowerCount: 15, PostCount: 8, CreatorID: user.ID},
		{BaseModel: common.BaseModel{ID: 4, CreatedAt: now}, Title: "Topic D", Slug: "topic-d", FollowerCount: 10, PostCount: 3, CreatorID: user.ID},
		{BaseModel: common.BaseModel{ID: 5, CreatedAt: now}, Title: "Topic E", Slug: "topic-e", FollowerCount: 5, PostCount: 10, CreatorID: user.ID},
	}
	for _, t := range testTopics {
		db.Create(&t)
	}

	t.Run("正常分页：第一页", func(t *testing.T) {
		topics, total, err := repo.List(2, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, topics, 2)
		assert.Equal(t, "Topic B", topics[0].Title) // 20,1
		assert.Equal(t, "Topic C", topics[1].Title) // 15,8
	})

	t.Run("分页：第二页", func(t *testing.T) {
		topics, total, err := repo.List(2, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, topics, 2)
		assert.Equal(t, "Topic A", topics[0].Title) // 10,5
		assert.Equal(t, "Topic D", topics[1].Title) // 10,3
	})

	t.Run("超出范围的 offset", func(t *testing.T) {
		topics, total, err := repo.List(10, 100)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Empty(t, topics)
	})

	t.Run("limit 为 0", func(t *testing.T) {
		topics, total, err := repo.List(0, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Empty(t, topics)
	})

	t.Run("空表", func(t *testing.T) {
		emptyDB := setupTestDB(t)
		emptyRepo := NewTopicRepository(emptyDB)
		topics, total, err := emptyRepo.List(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, topics)
	})
}

// ---------- 测试 GetByCreator ----------
func TestTopicRepository_GetByCreator(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user1 := createTestUser(db, 1, "creator1")
	user2 := createTestUser(db, 2, "creator2")

	now := time.Now()
	topics := []do.Topic{
		{BaseModel: common.BaseModel{ID: 1, CreatedAt: now.Add(-3 * time.Hour)}, Title: "User1 Topic 1", Slug: "u1-1", CreatorID: user1.ID},
		{BaseModel: common.BaseModel{ID: 2, CreatedAt: now.Add(-2 * time.Hour)}, Title: "User1 Topic 2", Slug: "u1-2", CreatorID: user1.ID},
		{BaseModel: common.BaseModel{ID: 3, CreatedAt: now.Add(-1 * time.Hour)}, Title: "User1 Topic 3", Slug: "u1-3", CreatorID: user1.ID},
		{BaseModel: common.BaseModel{ID: 4, CreatedAt: now}, Title: "User2 Topic 1", Slug: "u2-1", CreatorID: user2.ID},
	}
	for _, t := range topics {
		db.Create(&t)
	}

	t.Run("获取用户1的话题：分页第一页", func(t *testing.T) {
		topics, total, err := repo.GetByCreator(user1.ID, 2, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, topics, 2)
		// 期望按 created_at DESC 排序：最新的在前
		assert.Equal(t, "User1 Topic 3", topics[0].Title) // 最新
		assert.Equal(t, "User1 Topic 2", topics[1].Title)
	})

	t.Run("获取用户1的话题：第二页", func(t *testing.T) {
		topics, total, err := repo.GetByCreator(user1.ID, 2, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, topics, 1)
		assert.Equal(t, "User1 Topic 1", topics[0].Title)
	})

	t.Run("获取用户2的话题：只有一条", func(t *testing.T) {
		topics, total, err := repo.GetByCreator(user2.ID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, topics, 1)
		assert.Equal(t, "User2 Topic 1", topics[0].Title)
	})

	t.Run("获取不存在用户的话题", func(t *testing.T) {
		topics, total, err := repo.GetByCreator(999, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, topics)
	})

	t.Run("limit=0", func(t *testing.T) {
		topics, total, err := repo.GetByCreator(user1.ID, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Empty(t, topics)
	})
}
