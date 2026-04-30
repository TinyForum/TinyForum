package wire

import (
	"fmt"
	"time"

	"tiny-forum/config"
	adminInit "tiny-forum/init"

	// handler

	// repository

	// service

	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"

	_ "tiny-forum/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Private.Database.Host,
		cfg.Private.Database.User,
		cfg.Private.Database.Password,
		cfg.Private.Database.DBName,
		cfg.Private.Database.Port,
		cfg.Private.Database.SSLMode,
		cfg.Private.Database.TimeZone,
	)

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
		&model.RefreshToken{},
		&model.User{},                 // 用户
		&model.Follow{},               // 关注
		&model.Tag{},                  // 标签
		&model.Post{},                 // 帖子
		&model.Comment{},              // 评论
		&model.Like{},                 // 点赞
		&model.Notification{},         // 通知
		&model.SignIn{},               // 登录
		&model.Report{},               // 举报
		&model.Board{},                // 板块
		&model.Moderator{},            // 管理员
		&model.BoardBan{},             // 禁言
		&model.ModeratorLog{},         // 管理员日志
		&model.Question{},             // 问题
		&model.AnswerVote{},           // 回答投票
		&model.TimelineEvent{},        // 时间线事件
		&model.UserTimeline{},         // 用户时间线
		&model.TimelineSubscription{}, // 时间线订阅
		&model.Topic{},                // 主题
		&model.TopicPost{},            // 主题帖子
		&model.TopicFollow{},          // 主题关注
		&model.Announcement{},         // 公告
		&model.ModeratorApplication{}, // 版主申请
		&model.Moderator{},            // 版主
		&model.Vote{},
		// 审计
		&model.ContentAuditTask{}, // 内容审核任务
		&model.AuditLog{},         // 审计日志
		&model.UserRiskRecord{},   // 用户风险记录
		&model.Attachment{},       // 附件
		&model.IPRiskRecord{},     // IP风险记录
		&model.UserRiskRecord{},   // 用户风险记录
		&model.BlockedIP{},        // 被封禁IP
		&model.Violation{},        // 违规
		&model.Favorite{},         // 收藏
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
