package wire

import (
	"fmt"
	"time"

	"tiny-forum/config"
	adminInit "tiny-forum/init"
	"tiny-forum/internal/model/po"

	"tiny-forum/pkg/logger"

	_ "tiny-forum/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(&cfg.Private.Database)
	fmt.Printf("DSN config: %v\n", dsn)
	fmt.Printf("DBname: %v\n", cfg.Private.Database.DBName)
	logLevel := gormlogger.Silent
	if cfg.Basic.Server.Mode == "debug" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate all models
	if err := db.AutoMigrate(
		// 用户
		&po.RefreshToken{},
		&po.User{},                 // 用户
		&po.Follow{},               // 关注
		&po.Tag{},                  // 标签
		&po.Post{},                 // 帖子
		&po.Comment{},              // 评论
		&po.Like{},                 // 点赞
		&po.Notification{},         // 通知
		&po.SignIn{},               // 登录
		&po.Report{},               // 举报
		&po.Board{},                // 板块
		&po.Moderator{},            // 管理员
		&po.BoardBan{},             // 禁言
		&po.ModeratorLog{},         // 管理员日志
		&po.Question{},             // 问题
		&po.AnswerVote{},           // 回答投票
		&po.TimelineEvent{},        // 时间线事件
		&po.UserTimeline{},         // 用户时间线
		&po.TimelineSubscription{}, // 时间线订阅
		&po.Topic{},                // 主题
		&po.TopicPost{},            // 主题帖子
		&po.TopicFollow{},          // 主题关注
		&po.Announcement{},         // 公告
		&po.ModeratorApplication{}, // 版主申请
		&po.Moderator{},            // 版主
		&po.Vote{},
		// 审计
		&po.ContentAuditTask{}, // 内容审核任务
		&po.AuditLog{},         // 审计日志
		&po.UserRiskRecord{},   // 用户风险记录
		&po.Attachment{},       // 附件
		&po.IPRiskRecord{},     // IP风险记录
		&po.UserRiskRecord{},   // 用户风险记录
		&po.BlockedIP{},        // 被封禁IP
		&po.Violation{},        // 违规
		&po.Favorite{},         // 收藏
	); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	// 创建超级管理员
	if err := adminInit.CreateSuperAdmin(db, &cfg.Private.Admin); err != nil {
		logger.Warnf("创建超级管理员失败: %v", err)
	}

	// 核心配置：避免打爆 PostgreSQL
	sqlDB.SetMaxOpenConns(80)                 // 最大打开连接数（PG 默认 max_connections=100）
	sqlDB.SetMaxIdleConns(20)                 // 空闲连接池大小
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // 空闲连接超时

	logger.Info("Database connected and migrated successfully")
	return db, nil
}
