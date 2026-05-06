package init

import (
	"fmt"

	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminConfig 管理员配置

// DefaultAdminConfig 默认管理员配置
var DefaultAdminUserConfig = config.AdminUserConfig{
	Email:    "admin@test.com",
	Username: "admin",
	Password: "password",
	Role:     "super_admin", // 超级管理员
}

var DefaultSystemUserConfig = config.SystemUserConfig{
	Email:    "system@test.com",
	Username: "system",
	Password: "password",
	Role:     "system_maintainer", // 系统维护者
}

// createUserIfNotExists 公共的创建用户方法
// 如果用户已存在（通过 email 或 username 判断），则跳过创建并返回 nil
// 否则创建新用户，支持自定义额外字段（如 Score, Avatar 等）
func createUserIfNotExists(db *gorm.DB, email, username, password string, role do.UserRole, opts ...func(*do.User)) error {
	// 检查是否已存在
	var existingUser do.User
	result := db.Where("email = ? OR username = ?", email, username).First(&existingUser)

	if result.Error == nil {
		logger.Infof("用户已存在 (ID: %d, Email: %s, Role: %s)", existingUser.ID, existingUser.Email, existingUser.Role)
		return nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询用户失败: %w", result.Error)
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := &do.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		IsActive:  true,
		IsBlocked: false,
		Score:     0, // 默认积分，可由 opts 覆盖
	}

	// 默认头像根据用户名生成
	user.Avatar = fmt.Sprintf("https://api.dicebear.com/8.x/lorelei/svg?seed=%s", username)

	// 应用额外的配置选项
	for _, opt := range opts {
		opt(user)
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	logger.Infof("用户创建成功！ ID: %d, Email: %s, Role: %s", user.ID, user.Email, user.Role)
	return nil
}

// CreateSuperAdmin 创建超级管理员（复用公共方法）
func CreateSuperAdmin(db *gorm.DB, config *config.AdminUserConfig) error {
	if config == nil {
		config = &DefaultAdminUserConfig
	}
	// 为超级管理员设置高积分和特定角色
	return createUserIfNotExists(db, config.Email, config.Username, config.Password, config.Role,
		func(u *do.User) {
			u.Score = 10000 // 超级管理员初始积分
		},
	)
}

// CreateSystemUser 创建系统维护者（复用公共方法）
func CreateSystemUser(db *gorm.DB, config *config.SystemUserConfig) error {
	if config == nil {
		config = &DefaultSystemUserConfig
	}
	// 系统维护者不需要特殊积分，使用默认 0
	return createUserIfNotExists(db, config.Email, config.Username, config.Password, config.Role)
}
