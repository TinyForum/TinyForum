package topic

// go test -v ./internal/repository -run TestTopicRepository_List
import (
	"testing"
	"time"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// go test -v ./internal/repository/topic

// ---------- 测试 Follow ----------
func TestTopicRepository_Follow(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	// 创建测试用户
	user1 := createTestUser(db, 1, "follower1")
	user2 := createTestUser(db, 2, "follower2")

	// 创建测试话题
	topic := &do.Topic{
		Title:     "Follow Test Topic",
		Slug:      "follow-test-topic",
		CreatorID: user1.ID, // 创建者也是用户1
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	t.Run("首次关注成功", func(t *testing.T) {
		follow := &do.TopicFollow{
			UserID:  user2.ID,
			TopicID: topic.ID,
		}
		err := repo.Follow(follow)
		assert.NoError(t, err)
		assert.NotZero(t, follow.ID) // 应该被赋ID

		// 验证数据库中确实存在
		var count int64
		db.Model(&do.TopicFollow{}).Where("user_id = ? AND topic_id = ?", user2.ID, topic.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("重复关注 - 已存在则直接返回nil", func(t *testing.T) {
		// 已存在一条关注记录（上面创建的）
		follow := &do.TopicFollow{
			UserID:  user2.ID,
			TopicID: topic.ID,
		}
		err := repo.Follow(follow)
		assert.NoError(t, err)
		// 验证不会重复创建，ID仍为0（因为未创建新记录）
		assert.Equal(t, uint(0), follow.ID)

		// 验证数据库中只有一条记录
		var count int64
		db.Model(&do.TopicFollow{}).Where("user_id = ? AND topic_id = ?", user2.ID, topic.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("关注不存在的话题 - 外键约束失败", func(t *testing.T) {
		follow := &do.TopicFollow{
			UserID:  user2.ID,
			TopicID: 99999, // 不存在
		}
		err := repo.Follow(follow)
		assert.Error(t, err) // 应该返回外键约束错误
	})

	t.Run("关注不存在的用户 - 外键约束失败", func(t *testing.T) {
		follow := &do.TopicFollow{
			UserID:  99999,
			TopicID: topic.ID,
		}
		err := repo.Follow(follow)
		assert.Error(t, err)
	})
}

// ---------- 测试 Unfollow ----------
func TestTopicRepository_Unfollow(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user1 := createTestUser(db, 1, "user1")
	user2 := createTestUser(db, 2, "user2")
	topic := &do.Topic{
		Title:     "Unfollow Test",
		Slug:      "unfollow-test",
		CreatorID: user1.ID,
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	// 先创建一条关注记录
	follow := &do.TopicFollow{
		UserID:  user2.ID,
		TopicID: topic.ID,
	}
	err = db.Create(follow).Error
	require.NoError(t, err)

	t.Run("取消关注成功", func(t *testing.T) {
		err := repo.Unfollow(user2.ID, topic.ID)
		assert.NoError(t, err)

		// 验证记录已被软删除（DeletedAt 非空）
		var follow do.TopicFollow
		err = db.Unscoped().Where("user_id = ? AND topic_id = ?", user2.ID, topic.ID).First(&follow).Error
		assert.NoError(t, err)
		assert.NotNil(t, follow.DeletedAt.Time) // 或者 assert.False(t, follow.DeletedAt.Valid)
	})

	t.Run("取消不存在的关注 - 不会报错", func(t *testing.T) {
		err := repo.Unfollow(user2.ID, topic.ID) // 已经删除了
		assert.NoError(t, err)                   // GORM 删除不存在的记录影响0行，不返回错误
	})

	t.Run("取消关注时用户或话题不存在 - 不影响", func(t *testing.T) {
		err := repo.Unfollow(99999, topic.ID)
		assert.NoError(t, err)
	})
}

// ---------- 测试 IsFollowing ----------
func TestTopicRepository_IsFollowing(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	user1 := createTestUser(db, 1, "user1")
	user2 := createTestUser(db, 2, "user2")
	topic := &do.Topic{
		Title:     "IsFollowing Test",
		Slug:      "isfollowing-test",
		CreatorID: user1.ID,
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	// 创建一条关注记录
	follow := &do.TopicFollow{
		UserID:  user2.ID,
		TopicID: topic.ID,
	}
	err = db.Create(follow).Error
	require.NoError(t, err)

	t.Run("用户正在关注", func(t *testing.T) {
		following, err := repo.IsFollowing(user2.ID, topic.ID)
		assert.NoError(t, err)
		assert.True(t, following)
	})

	t.Run("用户未关注", func(t *testing.T) {
		following, err := repo.IsFollowing(user1.ID, topic.ID) // user1 未关注该话题（只是创建者）
		assert.NoError(t, err)
		assert.False(t, following)
	})

	t.Run("话题不存在 - 返回false", func(t *testing.T) {
		following, err := repo.IsFollowing(user2.ID, 99999)
		assert.NoError(t, err) // 查询不会报错，count为0
		assert.False(t, following)
	})

	t.Run("用户不存在 - 返回false", func(t *testing.T) {
		following, err := repo.IsFollowing(99999, topic.ID)
		assert.NoError(t, err)
		assert.False(t, following)
	})
}

// ---------- 测试 GetFollowers ----------
func TestTopicRepository_GetFollowers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTopicRepository(db)

	// 创建多个用户
	users := make([]*do.User, 5)
	for i := 1; i <= 5; i++ {
		users[i-1] = createTestUser(db, uint(i), "follower"+string(rune('a'+i-1)))
	}

	// 创建话题
	topic := &do.Topic{
		Title:     "GetFollowers Test",
		Slug:      "getfollowers-test",
		CreatorID: users[0].ID, // 用户1是创建者，但不一定关注
	}
	err := db.Create(topic).Error
	require.NoError(t, err)

	// 创建关注记录：用户2-5关注该话题，时间不同
	now := time.Now()
	for i := 1; i < 5; i++ {
		follow := &do.TopicFollow{
			BaseModel: common.BaseModel{CreatedAt: now.Add(time.Duration(i) * time.Second)}, // 不同时间
			UserID:    users[i].ID,
			TopicID:   topic.ID,
		}
		err = db.Create(follow).Error
		require.NoError(t, err)
	}

	t.Run("获取所有关注者 - 第一页", func(t *testing.T) {
		follows, total, err := repo.GetFollowers(topic.ID, 2, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, follows, 2)

		// 排序按 created_at DESC，最新的在前
		// 用户5最新，用户4次之
		assert.Equal(t, users[4].ID, follows[0].UserID) // 用户5
		assert.Equal(t, users[3].ID, follows[1].UserID) // 用户4
	})

	t.Run("获取所有关注者 - 第二页", func(t *testing.T) {
		follows, total, err := repo.GetFollowers(topic.ID, 2, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, follows, 2)

		assert.Equal(t, users[2].ID, follows[0].UserID) // 用户3
		assert.Equal(t, users[1].ID, follows[1].UserID) // 用户2
	})

	t.Run("超出范围offset", func(t *testing.T) {
		follows, total, err := repo.GetFollowers(topic.ID, 10, 100)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Empty(t, follows)
	})

	t.Run("limit=0", func(t *testing.T) {
		follows, total, err := repo.GetFollowers(topic.ID, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Empty(t, follows)
	})

	t.Run("话题没有关注者", func(t *testing.T) {
		// 创建一个新话题，不添加关注
		newTopic := &do.Topic{
			Title:     "Empty Topic",
			Slug:      "empty-topic",
			CreatorID: users[0].ID,
		}
		err := db.Create(newTopic).Error
		require.NoError(t, err)

		follows, total, err := repo.GetFollowers(newTopic.ID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, follows)
	})

	t.Run("验证预加载 User", func(t *testing.T) {
		follows, _, err := repo.GetFollowers(topic.ID, 1, 0)
		assert.NoError(t, err)
		assert.Len(t, follows, 1)
		// 验证 User 已被预加载
		assert.NotNil(t, follows[0].User)
		assert.Equal(t, users[4].ID, follows[0].User.ID)
		assert.Equal(t, users[4].Username, follows[0].User.Username)
	})
}
