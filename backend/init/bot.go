package initdata

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// InitDefaultBot 在数据库迁移后调用，确保系统默认机器人存在
func InitDefaultBot(db *gorm.DB) error {
	var bot do.Bot
	err := db.First(&bot, do.SystemBotID).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	defaultBot := do.DefaultSystemBot

	// 初始化所有 JSON 字段为非 nil
	if defaultBot.Screenshots == nil {
		defaultBot.Screenshots = []string{}
	}
	if defaultBot.Tags == nil {
		defaultBot.Tags = []string{}
	}
	if defaultBot.EnvVars == nil {
		defaultBot.EnvVars = make(map[string]string)
	}
	if defaultBot.Permissions == nil {
		defaultBot.Permissions = []do.BotPermission{}
	}
	if defaultBot.ConfigSchema == nil {
		defaultBot.ConfigSchema = []do.BotConfigField{}
	}
	if defaultBot.ConfigValues == nil {
		defaultBot.ConfigValues = make(map[string]any)
	}
	// ResourceLimit 可以为 nil，但 JSON 序列化会变成 null，PostgreSQL 可能接受
	// 如有必要，可赋值空对象，但建议保持 nil 并观察

	return db.Create(defaultBot).Error
}
