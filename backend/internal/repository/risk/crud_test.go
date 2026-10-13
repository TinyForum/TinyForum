package risk

import (
	"testing"

	"tiny-forum/internal/model/do"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupRiskDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&do.User{}, &do.ContentAuditTask{}))
	return db
}

// TestGetTaskByID 验证按 ID 获取审核任务，未命中返回 gorm.ErrRecordNotFound。
func TestGetTaskByID(t *testing.T) {
	db := setupRiskDB(t)
	repo := NewRiskRepository(db)

	task := &do.ContentAuditTask{
		TargetType:  do.AuditTargetPost,
		TargetID:    10,
		TriggerType: do.AuditTriggerSensitiveWord,
		Status:      do.ModerationStatusPending,
	}
	require.NoError(t, db.Create(task).Error)

	got, err := repo.GetTaskByID(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.ID, got.ID)
	assert.Equal(t, do.AuditTargetPost, got.TargetType)
	assert.EqualValues(t, 10, got.TargetID)

	_, err = repo.GetTaskByID(9999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

// TestListResolvedTasks 仅返回已处理（approved/rejected）的任务，按审核时间倒序。
func TestListResolvedTasks(t *testing.T) {
	db := setupRiskDB(t)
	repo := NewRiskRepository(db)

	tasks := []*do.ContentAuditTask{
		{TargetType: do.AuditTargetPost, TargetID: 1, TriggerType: do.AuditTriggerManual, Status: do.ModerationStatusPending},
		{TargetType: do.AuditTargetPost, TargetID: 2, TriggerType: do.AuditTriggerManual, Status: do.ModerationStatusApproved},
		{TargetType: do.AuditTargetComment, TargetID: 3, TriggerType: do.AuditTriggerManual, Status: do.ModerationStatusRejected},
	}
	for _, task := range tasks {
		require.NoError(t, db.Create(task).Error)
	}

	got, total, err := repo.ListResolvedTasks(20, 0)
	require.NoError(t, err)
	assert.EqualValues(t, 2, total)
	assert.Len(t, got, 2)
	for _, task := range got {
		assert.NotEqual(t, do.ModerationStatusPending, task.Status)
	}
}
