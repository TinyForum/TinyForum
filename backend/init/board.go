package initdata

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// InitDefaultBot 确保系统默认机器人存在，并返回该机器人实例
func InitDefaultBorad(db *gorm.DB) (*do.Board, error) {
	var board do.Board
	err := db.First(&board, do.SystemBoardID).Error
	if err == nil {
		return &board, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 记录不存在，创建默认机器人
	defaultWorldBoard := do.DefaultWorldBoard

	// 初始化所有 JSON 字段为非 nil

	// 创建记录（注意 Create 传入指针）
	if err := db.Create(&defaultWorldBoard).Error; err != nil {
		return nil, err
	}
	return defaultWorldBoard, nil
}
