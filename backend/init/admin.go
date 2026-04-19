package init

import (
	"fmt"

	"tiny-forum/config"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminConfig 管理员配置

// DefaultAdminConfig 默认管理员配置
var DefaultAdminConfig = config.AdminConfig{
	Email:    "admin@test.com",
	Username: "admin",
	Password: "password",
	Role:     "super_admin", // 或 "super_admin"，根据你的角色定义
}

// CreateSuperAdmin 创建超级管理员
func CreateSuperAdmin(db *gorm.DB, config *config.AdminConfig) error {
	if config == nil {
		config = &DefaultAdminConfig
	}

	// 检查是否已存在
	var existingUser model.User
	result := db.Where("email = ? OR username = ?", config.Email, config.Username).First(&existingUser)

	if result.Error == nil {
		logger.Infof("超级管理员已存在 (ID: %d, Email: %s)", existingUser.ID, existingUser.Email)
		return nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询用户失败: %w", result.Error)
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建管理员
	admin := &model.User{
		Username:  config.Username,
		Email:     config.Email,
		Password:  string(hashedPassword),
		Role:      config.Role,
		IsActive:  true,
		IsBlocked: false,
		Score:     10000,
	}

	// 设置默认头像
	admin.Avatar = fmt.Sprintf("https://api.dicebear.com/8.x/lorelei/svg?seed=%s", config.Username)

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("创建管理员失败: %w", err)
	}

	logger.Infof("超级管理员创建成功！\n  ID: %d\n  Email: %s\n  密码: %s", admin.ID, admin.Email, config.Password)
	return nil
}
