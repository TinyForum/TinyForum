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
		&model.User{},
		&model.Follow{},
		&model.Tag{},
		&model.Post{},
		&model.Comment{},
		&model.Like{},
		&model.Notification{},
		&model.SignIn{},
		&model.Report{},
		&model.Board{},
		&model.Moderator{},
		&model.BoardBan{},
		&model.ModeratorLog{},
		&model.Question{},
		&model.AnswerVote{},
		&model.TimelineEvent{},
		&model.UserTimeline{},
		&model.TimelineSubscription{},
		&model.Topic{},
		&model.TopicPost{},
		&model.TopicFollow{},
		&model.Announcement{},
		&model.ModeratorApplication{},
		&model.Moderator{},
		&model.Vote{},
		// 审计
		&model.ContentAuditTask{},
		&model.AuditLog{},
		&model.UserRiskRecord{},
		&model.Attachment{},
		&model.IPRiskRecord{},
		&model.UserRiskRecord{},
		&model.BlockedIP{},
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
