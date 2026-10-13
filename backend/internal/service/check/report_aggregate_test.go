package check

import (
	"testing"

	"tiny-forum/internal/model/do"
	riskrepo "tiny-forum/internal/repository/risk"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCheckDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&do.User{}, &do.Report{}, &do.ContentAuditTask{}))
	return db
}

// TestHandleReportAggregate_TriggersAtThreshold 举报数达到阈值时应创建审核任务且保持幂等。
func TestHandleReportAggregate_TriggersAtThreshold(t *testing.T) {
	db := setupCheckDB(t)
	svc := NewContentCheckService(riskrepo.NewRiskRepository(db), nil)

	// 未达阈值：不触发
	for i := 0; i < do.ReportAggregateThreshold-1; i++ {
		require.NoError(t, db.Create(&do.Report{
			ReporterID: uint(i + 1),
			TargetID:   10,
			TargetType: do.ReportTargetTypePost,
			Reason:     "spam",
			Status:     do.ReportPending,
		}).Error)
	}
	triggered, err := svc.HandleReportAggregate(do.AuditTargetPost, 10)
	require.NoError(t, err)
	assert.False(t, triggered)

	var count int64
	require.NoError(t, db.Model(&do.ContentAuditTask{}).Count(&count).Error)
	assert.EqualValues(t, 0, count)

	// 达到阈值：触发并创建 pending 任务
	require.NoError(t, db.Create(&do.Report{
		ReporterID: 99,
		TargetID:   10,
		TargetType: do.ReportTargetTypePost,
		Reason:     "spam",
		Status:     do.ReportPending,
	}).Error)

	triggered, err = svc.HandleReportAggregate(do.AuditTargetPost, 10)
	require.NoError(t, err)
	assert.True(t, triggered)

	require.NoError(t, db.Model(&do.ContentAuditTask{}).
		Where("target_type = ? AND target_id = ? AND status = ?", do.AuditTargetPost, 10, do.ModerationStatusPending).
		Count(&count).Error)
	assert.EqualValues(t, 1, count)

	// 再次触发：幂等，不重复创建
	triggered, err = svc.HandleReportAggregate(do.AuditTargetPost, 10)
	require.NoError(t, err)
	assert.True(t, triggered)
	require.NoError(t, db.Model(&do.ContentAuditTask{}).Count(&count).Error)
	assert.EqualValues(t, 1, count)
}
