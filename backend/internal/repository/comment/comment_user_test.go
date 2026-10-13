package comment

import (
	"context"
	"testing"

	"tiny-forum/internal/model/do"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCommentUserDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&do.User{}, &do.Reply{}, &do.Comment{}))
	return db
}

func createCommentByAuthor(t *testing.T, db *gorm.DB, authorID, worksID uint, content string) do.Comment {
	t.Helper()
	reply := do.Reply{Content: content, AuthorID: authorID, TargetType: do.TargetReplyTypePost, TargetID: worksID}
	require.NoError(t, db.Create(&reply).Error)
	comment := do.Comment{WorksID: worksID, ReplyID: reply.ID}
	require.NoError(t, db.Create(&comment).Error)
	return comment
}

// TestListByAuthor_ReturnsOnlyOwnComments 验证按作者过滤，且分页总数正确。
func TestListByAuthor_ReturnsOnlyOwnComments(t *testing.T) {
	db := setupCommentUserDB(t)
	repo := &commentRepository{db: db}

	createCommentByAuthor(t, db, 1, 10, "我的第一条评论")
	createCommentByAuthor(t, db, 1, 11, "我的第二条评论")
	createCommentByAuthor(t, db, 2, 10, "别人的评论")

	comments, total, err := repo.ListByAuthor(context.Background(), 1, 0, 20)
	require.NoError(t, err)
	assert.EqualValues(t, 2, total, "只应统计作者 1 的评论")
	require.Len(t, comments, 2)
	for _, c := range comments {
		require.NotNil(t, c.Reply)
		assert.EqualValues(t, 1, c.Reply.AuthorID)
	}

	other, otherTotal, err := repo.ListByAuthor(context.Background(), 2, 0, 20)
	require.NoError(t, err)
	assert.EqualValues(t, 1, otherTotal)
	require.Len(t, other, 1)
}

// TestListByAuthor_Pagination 验证分页偏移与截止正确。
func TestListByAuthor_Pagination(t *testing.T) {
	db := setupCommentUserDB(t)
	repo := &commentRepository{db: db}

	for i := 0; i < 5; i++ {
		createCommentByAuthor(t, db, 1, uint(10+i), "评论")
	}

	page1, total, err := repo.ListByAuthor(context.Background(), 1, 0, 2)
	require.NoError(t, err)
	assert.EqualValues(t, 5, total)
	assert.Len(t, page1, 2)

	page3, _, err := repo.ListByAuthor(context.Background(), 1, 4, 2)
	require.NoError(t, err)
	assert.Len(t, page3, 1, "最后一页应只剩 1 条")
}

// TestListByAuthor_ExcludesDeletedReply 验证软删除的回复不出现在列表中。
func TestListByAuthor_ExcludesDeletedReply(t *testing.T) {
	db := setupCommentUserDB(t)
	repo := &commentRepository{db: db}

	comment := createCommentByAuthor(t, db, 1, 10, "待删除")
	require.NoError(t, db.Delete(&do.Reply{}, comment.ReplyID).Error)

	comments, total, err := repo.ListByAuthor(context.Background(), 1, 0, 20)
	require.NoError(t, err)
	assert.EqualValues(t, 0, total)
	assert.Empty(t, comments)
}

// TestUpdateReplyStatusByCommentID 验证审核可联动更新回复显示状态。
func TestUpdateReplyStatusByCommentID(t *testing.T) {
	db := setupCommentUserDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	comment := createCommentByAuthor(t, db, 1, 10, "待审核评论")

	require.NoError(t, repo.UpdateReplyStatusByCommentID(ctx, comment.ID, do.ReplyStatusHidden))

	var reply do.Reply
	require.NoError(t, db.First(&reply, comment.ReplyID).Error)
	assert.Equal(t, do.ReplyStatusHidden, reply.Status)

	require.NoError(t, repo.UpdateReplyStatusByCommentID(ctx, comment.ID, do.ReplyStatusVisible))
	require.NoError(t, db.First(&reply, comment.ReplyID).Error)
	assert.Equal(t, do.ReplyStatusVisible, reply.Status)
}
