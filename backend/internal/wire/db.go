package wire

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	adminInit "tiny-forum/init"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/do"

	"tiny-forum/pkg/logger"

	_ "tiny-forum/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(&cfg.Postgres)
	if cfg.Postgres.Logger == nil {
		logger.Info("Postgres logger is not configured, using default settings")
		cfg.Postgres.Logger = &config.PostLogger{
			SlowThreshold:             "200ms",
			LogLevel:                  "slient",
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}
	}

	// 解析慢查询阈值
	slowThreshold, err := time.ParseDuration(cfg.Postgres.Logger.SlowThreshold)
	if err != nil {
		slowThreshold = 200 * time.Millisecond // 解析失败时的后备默认值
	}

	// 创建 GORM 日志配置
	gormLogConfig := gormlogger.Config{
		SlowThreshold:             slowThreshold,
		LogLevel:                  parseLogLevel(cfg.Postgres.Logger.LogLevel),
		IgnoreRecordNotFoundError: cfg.Postgres.Logger.IgnoreRecordNotFoundError,
		Colorful:                  cfg.Postgres.Logger.Colorful,
	}
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogConfig,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Auto migrate all models
	if err := db.AutoMigrate(
		// 用户
		&do.RefreshToken{},
		&do.User{},                 // 用户
		&do.Follow{},               // 关注
		&do.Tag{},                  // 标签
		&do.Post{},                 // 帖子
		&do.Comment{},              // 评论
		&do.Like{},                 // 点赞
		&do.Notification{},         // 通知
		&do.SignIn{},               // 登录
		&do.Report{},               // 举报
		&do.Board{},                // 板块
		&do.Moderator{},            // 管理员
		&do.BoardBan{},             // 禁言
		&do.ModeratorLog{},         // 管理员日志
		&do.Question{},             // 问题
		&do.AnswerVote{},           // 回答投票
		&do.TimelineEvent{},        // 时间线事件
		&do.UserTimeline{},         // 用户时间线
		&do.TimelineSubscription{}, // 时间线订阅
		&do.Topic{},                // 主题
		&do.TopicPost{},            // 主题帖子
		&do.TopicFollow{},          // 主题关注
		&do.Announcement{},         // 公告
		&do.ModeratorApplication{}, // 版主申请
		&do.Moderator{},            // 版主
		&do.Vote{},
		// 审计
		&do.ContentAuditTask{}, // 内容审核任务
		&do.AuditLog{},         // 审计日志
		&do.UserRiskRecord{},   // 用户风险记录
		&do.Attachment{},       // 附件
		&do.IPRiskRecord{},     // IP风险记录
		&do.UserRiskRecord{},   // 用户风险记录
		&do.BlockedIP{},        // 被封禁IP
		&do.Violation{},        // 违规
		&do.Favorite{},         // 收藏
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

func parseLogLevel(level string) gormlogger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Info // 默认 info
	}
}
