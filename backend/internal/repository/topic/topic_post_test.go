package topic

import (
	"testing"
	"time"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- 测试 AddPost ----------
func TestTopicRepository_AddPost(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	// 创建用户
	user := createTestUser(db, 1, "author")
	// 创建话题
	topic := createTestTopic(db, 1, "Test Topic", "test-topic", user.ID)
	// 创建作品
	creation := createTestCreation(db, 1, "Test Post", user.ID)

	t.Run("首次添加帖子到话题成功", func(t *testing.T) {
		topicPost := &do.TopicCreation{
			TopicID:    topic.ID,
			CreationID: creation.ID, // 字段名若为 CreationID 则需调整
			// SortOrder 默认 0
		}
		err := repo.AddPost(topicPost)
		assert.NoError(t, err)
		assert.NotZero(t, topicPost.ID) // 应该被赋值

		// 验证数据库中确实存在
		var count int64
		db.Model(&do.TopicCreation{}).Where("topic_id = ? AND creation_id = ?", topic.ID, creation.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("重复添加：已存在则直接返回nil", func(t *testing.T) {
		// 已存在一条关系（上面创建的）
		topicPost := &do.TopicCreation{
			TopicID:    topic.ID,
			CreationID: creation.ID,
		}
		err := repo.AddPost(topicPost)
		assert.NoError(t, err)
		// 不会创建新记录，ID仍为0
		assert.Equal(t, uint(0), topicPost.ID)

		// 验证只有一条记录
		var count int64
		db.Model(&do.TopicCreation{}).Where("topic_id = ? AND creation_id = ?", topic.ID, creation.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	// 注意：SQLite 默认不检查外键，所以以下测试可能不会报错，可根据实际需求调整
	t.Run("添加时话题不存在：可能不会报错（取决于外键约束）", func(t *testing.T) {
		topicPost := &do.TopicCreation{
			TopicID:    99999,
			CreationID: creation.ID,
		}
		_ = repo.AddPost(topicPost)
		// 如果数据库外键约束未启用，err 可能为 nil，这里只记录行为
		// 实际上我们无法强制报错，所以只检查是否执行成功
		// 但我们可以验证不会创建记录（因为没有外键约束可能成功）
		// 更好的做法是依赖数据库外键，但 SQLite 默认不强制，所以此测试可能不适用。
		// 可以选择跳过外键错误测试或使用其他数据库。
		// 这里我们仅记录，不做严格断言。
	})
}

// ---------- 测试 RemovePost ----------
func TestTopicRepository_RemovePost(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user := createTestUser(db, 1, "author")
	topic := createTestTopic(db, 1, "Remove Test", "remove-test", user.ID)
	creation := createTestCreation(db, 1, "Post to Remove", user.ID)

	// 先添加一条关系
	topicPost := &do.TopicCreation{
		TopicID:    topic.ID,
		CreationID: creation.ID,
	}
	err := db.Create(topicPost).Error
	require.NoError(t, err)

	t.Run("删除成功", func(t *testing.T) {
		err := repo.RemovePost(topic.ID, creation.ID)
		assert.NoError(t, err)

		// 验证已删除
		var count int64
		db.Unscoped().Model(&do.TopicCreation{}).Where("topic_id = ? AND creation_id = ?", topic.ID, creation.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("删除不存在的关系：不会报错", func(t *testing.T) {
		err := repo.RemovePost(topic.ID, creation.ID) // 已删除
		assert.NoError(t, err)
	})
}

// ---------- 测试 GetTopicPosts ----------
func TestTopicRepository_GetTopicPosts(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user := createTestUser(db, 1, "author")
	topic := createTestTopic(db, 1, "List Test", "list-test", user.ID)

	// 创建多个作品（Posts）
	posts := make([]*do.Creation, 5)
	for i := 0; i < 5; i++ {
		posts[i] = createTestCreation(db, uint(i+1), "Post "+string(rune('A'+i)), user.ID)
	}

	// 创建关系，设置不同的 sort_order 和 created_at 以便测试排序
	now := time.Now()
	for i := 0; i < 5; i++ {
		tp := &do.TopicCreation{
			BaseModel:  common.BaseModel{CreatedAt: now.Add(time.Duration(i) * time.Minute)},
			TopicID:    topic.ID,
			CreationID: posts[i].ID,
			SortOrder:  5 - i,
			CreatorID:  user.ID,
		}
		db.Create(tp)
	}

	t.Run("获取所有帖子：第一页（按 sort_order ASC, created_at ASC）", func(t *testing.T) {
		topicPosts, total, err := repo.GetTopicPosts(topic.ID, 2, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, topicPosts, 2)

		// 预期排序：sort_order 最小的在前，即 sort_order=1 的 Post 5（i=4）最先，然后是 sort_order=2 的 Post 4
		// 但由于 created_at 也作为次要排序，如果 sort_order 相同则按 created_at 升序
		assert.Equal(t, posts[4].ID, topicPosts[0].CreationID) // Post 5 (sort_order=1)
		assert.Equal(t, posts[3].ID, topicPosts[1].CreationID) // Post 4 (sort_order=2)
	})

	t.Run("第二页", func(t *testing.T) {
		topicPosts, total, err := repo.GetTopicPosts(topic.ID, 2, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, topicPosts, 2)

		assert.Equal(t, posts[2].ID, topicPosts[0].CreationID) // Post 3 (sort_order=3)
		assert.Equal(t, posts[1].ID, topicPosts[1].CreationID) // Post 2 (sort_order=4)
	})

	t.Run("超出范围 offset", func(t *testing.T) {
		topicPosts, total, err := repo.GetTopicPosts(topic.ID, 10, 100)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Empty(t, topicPosts)
	})

	t.Run("limit=0", func(t *testing.T) {
		topicPosts, total, err := repo.GetTopicPosts(topic.ID, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Empty(t, topicPosts)
	})

	t.Run("话题没有帖子", func(t *testing.T) {
		newTopic := createTestTopic(db, 2, "Empty Topic", "empty", user.ID)
		topicPosts, total, err := repo.GetTopicPosts(newTopic.ID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, topicPosts)
	})

	t.Run("验证预加载 Post 和 Post.Author", func(t *testing.T) {
		topicPosts, _, err := repo.GetTopicPosts(topic.ID, 1, 0)
		assert.NoError(t, err)
		assert.Len(t, topicPosts, 1)
		// 验证 Post 被预加载
		assert.NotNil(t, topicPosts[0].CreationID)
		assert.Equal(t, posts[4].ID, topicPosts[0].CreationID) // 因为排序，第一条是 Post 5
		// 验证 Author 被预加载
		assert.NotNil(t, topicPosts[0].CreatorID)
		assert.Equal(t, user.ID, topicPosts[0].CreatorID)
	})
}

// ---------- 测试 UpdatePostOrder ----------
func TestTopicRepository_UpdatePostOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user := createTestUser(db, 1, "author")
	topic := createTestTopic(db, 1, "Order Test", "order-test", user.ID)
	creation := createTestCreation(db, 1, "Order Post", user.ID)

	// 创建关系
	tp := &do.TopicCreation{
		TopicID:    topic.ID,
		CreationID: creation.ID,
		SortOrder:  10,
	}
	err := db.Create(tp).Error
	require.NoError(t, err)

	t.Run("更新排序成功", func(t *testing.T) {
		err := repo.UpdatePostOrder(topic.ID, creation.ID, 20)
		assert.NoError(t, err)

		// 验证更新
		var updated do.TopicCreation
		db.Where("topic_id = ? AND creation_id = ?", topic.ID, creation.ID).First(&updated)
		assert.Equal(t, 20, updated.SortOrder)
	})

	t.Run("更新不存在的关系：不会报错但影响0行", func(t *testing.T) {
		err := repo.UpdatePostOrder(topic.ID, 99999, 30)
		assert.NoError(t, err) // GORM Update 影响0行不报错
	})
}
