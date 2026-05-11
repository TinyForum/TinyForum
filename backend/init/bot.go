package initdata

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// InitDefaultBot 确保系统默认机器人存在，并返回该机器人实例
func InitDefaultBot(db *gorm.DB) (*do.Bot, error) {
	var bot do.Bot
	err := db.First(&bot, do.SystemBotID).Error
	if err == nil {
		return &bot, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 记录不存在，创建默认机器人
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

	// 创建记录（注意 Create 传入指针）
	if err := db.Create(&defaultBot).Error; err != nil {
		return nil, err
	}
	return defaultBot, nil
}
